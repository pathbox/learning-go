package main

import (
	"github.com/olivere/elastic"
	_ "log"
	"os"
	"strings"
	"time"
)

const (
	indexName = "films"
	mapping   = `
	{
		"settings":{
			"number_of_shards":1,
			"number_of_replicas":0
		},
		"mappings":{
			"_doc":{
				"properties":{
					"title":{
						"type":"keyword"
					},
					"genre":{
						"type":"keyword"
					},
					"year":{
						"type":"long"
					},
					"director":{
						"type":"keyword"
					}
				}
			}
		}
	}
	`
)

func main() {
	opts := []elastic.ClientOptionFunc{
		elastic.SetTraceLog(log.New(os.Stdout, "", 0)),
	}
	client, err := elastic.NewClient(opts...)
	if err != nil {
		panic(err)
	}

	err = createAndPopulateIndex(client)
	if err != nil {
		panic(err)
	}

	f := NewFinder()
	f = f.From(0).Size(100)
	f = f.Pretty(true)

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	res, err := f.Find(ctx, client)
	if err != nil {
		panic(err)
	}

	// Output results
	fmt.Printf("Searched through %d films\n", res.Total)
	fmt.Println()
	fmt.Println("Films found:")
	for i, film := range res.Films{
		prefix := "├"
		if i == len(res.Films)-1 {
			prefix = "└"
		}
		fmt.Printf("%s %s from %d\n", prefix, film.Title, film.Year)
	}
	fmt.Println()
	fmt.Println("Broken down by genre:")
	for year, genre := range res.YearsAndGenres {
		fmt.Printf("- %4d\n", year)
		for i, nc := range genre {
			prefix := "├"
			if i == len(genre)-1 {
				prefix = "└"
			}
			fmt.Printf("  %s%2d× %s\n", prefix, nc.Count, nc.Name)
		}
	}
}

// Film represents a movie with some properties.
type Film struct {
	Title string `json:"title"`
	Genre []string `json:"genre"`
	Year int `json:"year"`
	Director string `json:"director"`
}

type Finder struct {
	genre string
	year int
	from, size int
	sort []string
	pretty bool
}
// FinderResponse is the outcome of calling Finder.Find.
type FinderResponse struct {
	Total int64
	Films []*Film
	Genres map[string]int64
	YearsAndGenres map[int][]NameCount
}

type NameCount struct {
	Name string
	Count int64
}

// NewFinder creates a new finder for films.
// Use the funcs to set up filters and search properties,
// then call Find to execute.
func NewFinder() *Finder {
	return &Finder{}
}

func (f *Finder) Genre(genre string) *Finder {
	f.genre = genre
	return f
}

// Year filters the results by the specified year.
func (f *Finder) Year(year int) *Finder {
	f.year = year
	return f
}

// From specifies the start index for pagination.
func (f *Finder) From(from int) *Finder {
	f.from = from
	return f
}

// Size specifies the number of items to return in pagination.
func (f *Finder) Size(size int) *Finder {
	f.size = size
	return f
}

// Sort specifies one or more sort orders.
// Use a dash (-) to make the sort order descending.
// Example: "name" or "-year".
func (f *Finder) Sort(sort ...string) *Finder {
	if f.sort == nil {
		f.sort = make([]string, 0)
	}
	f.sort = append(f.sort, sort...)
	return f
}

// Pretty, when enabled, asks the server to return the
// response formatted and indented.
func (f *Finder) Pretty(pretty bool) *Finder {
	f.pretty = pretty
	return f
}

func (f *Finder) Find(ctx context.Context, client *elastic.Client) (FinderResponse,error) {
	var resp FinderResponse

	search := client.Search().Index(indexName).Type("_doc").Pretty(f.pretty)
	search = f.query(search)
	search = f.aggs(search)
	search = f.sorting(search)
	search = f.paginate(search)

	sr, err := search.Do(ctx)
	if err != nil {
		return resp, err
	}

	// Decode response
	films, err := f.decodeFilms(sr)
	if err != nil {
		return resp, err
	}
	resp.Films = films
	resp.Total = sr.Hits.TotalHits

	// Deserialize aggregations
	if agg, found := sr.Aggregations.Terms("all_genres"); found {
		resp.Genres = make(map[string]int64)
		for _, bucket := range agg.Buckets {
			resp.Genres[bucket.Key.(string)] = bucket.DocCount
		}
	}

	// Use the correct function on sr.Aggregations.XXX. It must match the
	// aggregation type specified at query time.
	// See https://github.com/olivere/elastic/blob/release-branch.v6/search_aggs.go
	// for all kinds of aggregation types.
	if agg, found := sr.Aggregations.Terms("years_and_genres"); found {
		resp.YearsAndGenres = make(map[int][]NameCount)
		for _, bucket := range agg.Buckets {
			// JSON doesn't have integer types: All numeric values are float64
			floatValue, ok := bucket.Key.(float64)
			if !ok {
				panic("expected a float64")
			}
			var (
				year          = int(floatValue)
				genresForYear []NameCount
			)
			// Iterate over the sub-aggregation
			if subAgg, found := bucket.Terms("genres_by_year"); found {
				for _, subBucket := range subAgg.Buckets {
					genresForYear = append(genresForYear, NameCount{
						Name:  subBucket.Key.(string),
						Count: subBucket.DocCount,
					})
				}
			}
			resp.YearsAndGenres[year] = genresForYear
		}
	}

	return resp, nil
}

func (f *Finder) query(service *elastic.SearchService) *elastic.SearchService {
	if f.genre =="" &&f.year == 0 {
		service = service.Query(elastic.NewMatchAllQuery())
		return service
	}

	q := elastic.NewBoolQuery()
	if f.genre != ""{
		q = q.Must(elastic.NewTermQuery("genre", f.genre))
	}
	if f.year > 0 {
		q = q.Must(elastic.NewTermQuery("year", f.year))
	}

	service = service.Query(q)
	return service
}

// aggs sets up the aggregations in the service.
func (f *Finder) aggs(service *elastic.SearchService) *elastic.SearchService {
	// Terms aggregation by genre
	agg := elastic.NewTermsAggregation().Field("genre")
	service = service.Aggregation("all_genres", agg)

	// Add a terms aggregation of Year, and add a sub-aggregation for Genre
	subAgg := elastic.NewTermsAggregation().Field("genre")
	agg = elastic.NewTermsAggregation().Field("year").
		SubAggregation("genres_by_year", subAgg)
	service = service.Aggregation("years_and_genres", agg)

	return service
}

// paginate sets up pagination in the service.
func (f *Finder) paginate(service *elastic.SearchService) *elastic.SearchService {
	if f.from > 0 {
		service = service.From(f.from)
	}
	if f.size > 0 {
		service = service.Size(f.size)
	}
	return service
}

// sorting applies sorting to the service.
func (f *Finder) sorting(service *elastic.SearchService) *elastic.SearchService {
	if len(f.sort) == 0 {
		// Sort by score by default
		service = service.Sort("_score", false)
		return service
	}

	// Sort by fields; prefix of "-" means: descending sort order.
	for _, s := range f.sort {
		s = strings.TrimSpace(s)

		var field string
		var asc bool

		if strings.HasPrefix(s, "-") {
			field = s[1:]
			asc = false
		} else {
			field = s
			asc = true
		}

		// Maybe check for permitted fields to sort

		service = service.Sort(field, asc)
	}
	return service
}

// decodeFilms takes a search result and deserializes the films.
func (f *Finder) decodeFilms(res *elastic.SearchResult) ([]*Film, error) {
	if res == nil || res.TotalHits() == 0 {
		return nil, nil
	}

	var films []*Film
	for _, hit := range res.Hits.Hits {
		film := new(Film)
		if err := json.Unmarshal(*hit.Source, film); err != nil {
			return nil, err
		}
		// TODO Add Score here, e.g.:
		// film.Score = *hit.Score
		films = append(films, film)
	}
	return films, nil
}

func createAndPopulateIndex(client *elastic.Client) error {
	ctx := context.Background()
	exists, err := client.IndexExists(indexName).Do(ctx)
	if err != nil {
		return err
	}
	if exists {
		_, err = client.DeleteIndex(indexName).Do(ctx)
		if err != nil {
			return err
		}
	}
	// Create index with mapping
	_, err = client.CreateIndex(indexName).Body(mapping).Do(ctx)
	if err != nil {
		return err
	}
	// Populate some films
	films := []Film{
		{Title: "The Shawshank Redemption", Genre: []string{"Crime", "Drama"}, Year: 1994, Director: "Frank Darabont"},
		{Title: "The Godfather", Genre: []string{"Crime", "Drama"}, Year: 1972, Director: "Francis Ford Coppola"},
		{Title: "The Godfather: Part II", Genre: []string{"Crime", "Drama"}, Year: 1974, Director: "Francis Ford Coppola"},
		{Title: "The Dark Knight", Genre: []string{"Action", "Crime", "Drama"}, Year: 2008, Director: "Christopher Nolan"},
		{Title: "12 Angry Men", Genre: []string{"Crime", "Drama"}, Year: 1957, Director: "Sidney Lumet"},
		{Title: "Schindler's List", Genre: []string{"Biography", "Drama", "History"}, Year: 1993, Director: "Steven Spielberg"},
		{Title: "The Lord of the Rings: The Return of the King", Genre: []string{"Adventure", "Drama", "Fantasy"}, Year: 2003, Director: "Peter Jackson"},
		{Title: "Pulp Fiction", Genre: []string{"Crime", "Drama"}, Year: 1994, Director: "Quentin Tarantino"},
		{Title: "Il buono, il brutto, il cattivo", Genre: []string{"Western"}, Year: 1966, Director: "Sergio Leone"},
		{Title: "Fight Club", Genre: []string{"Drama"}, Year: 1999, Director: "David Fincher"},
	}
	for _, film := range films {
		_, err = client.Index().
			Index(indexName).
			Type("_doc").
			BodyJson(film).
			Do(ctx)
		if err != nil {
			return err
		}
	}
	_, err = client.Flush(indexName).WaitIfOngoing(true).Do(ctx)
	return err
}











