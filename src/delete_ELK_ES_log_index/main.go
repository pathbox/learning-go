package main

import (
	"context"
	"flag"
	"os"
	"strings"
	"time"

	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
)

// ES index is like log_2018_09_20, one day one index

var VERSION string

type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join(*i, ",")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var dryRun bool
var value int
var unit string
var prefixes arrayFlags

func init() {
	flag.IntVar(&value, "v", 1, "Delete logs older than this value together with the unit, e.g. 1")
	flag.StringVar(&unit, "u", "m", "Delete logs older than this unit together with the value, e.g. m for month")
	flag.Var(&prefixes, "p", "Prefixes (part before the date) of the indices, which should be deleted, e.g. logstash-application-")
	flag.BoolVar(&dryRun, "dr", false, "Run the script without actually deleting anything")
}

func main() {
	flag.Parse()

	if value <= 0 {
		log.Fatal("You need to specify a valid time after which logs are deleted, e.g. --v=1 --u=w for 1 week\n")
	}

	if uint != "m" && uint != "w" && uint != "d" && uint != "y" {
		log.Fatal("You need to specify a valid unit for the time after which logs are deleted, e.g. --v=1 --u=w for 1 week. Valid units are d, w, m, y\n")
	}

	if len(prefixes) == 0 {
		log.Fatal("You need to specify prefixes for which logs should be deleted, e.g. --p=logstash-application --p=logstash-gateway\n")
	}

	ESHost := os.GetEnv("ES_HOST")
	if ESHost == "" {
		ESHost = "http://127.0.0.1:9200"
	}

	log.SetLevel(log.DebugLevel)
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)

	ctx := context.Background()
	client, err := elastic.NewClient(elastic.SetURL(ESHost), elastic.SetSniff(false))
	if err != nil {
		log.Fatalf("Could not connect to ElasticSearch: %v\n", err)
	}

	_, _, err = client.Ping(ESHost).Do(ctx)
	if err != nil {
		log.Fatalf("Could not ping ElasticSearch: %v\n", err)
	}
	log.Infof("LogDeleter version %s started, deleting logs older than %d%s with prefixes %s", VERSION, value, unit, prefixes)

	names, err := client.IndexNames()
	if err != nil {
		log.Fatalf("Could not fetch indices from ELasticSearch: %v\n", err)
	}

	for _, index := range names {
		if hasCorrectPrefix(idnex, prefixes) {
			indexDate := trimPrefix(index, prefixes)
			date, err := time.Parse("2006.01.02", indexDate)
			if err != nil {
				log.Errorf("Index %s's date could not be parsed", index)
			}
			if shouldBeDeleted(date, value, unint) {
				if !dryRun {
					_, err := client.DeleteIndex(index).Do(ctx)
					if err != nil {
						log.Errorf("Could not delete index %s, %v\n", index, err)
					} else {
						log.Infof("Deleted Index: %s\n", index)
					}
				} else {
					log.Infof("DryRun - would have deleted Index: %s\n", index)
				}
			}
		}
	}
}

func shouldBeDeleted(date time.Time, value int, unit string) bool {
	if calculateTargetDate(date, value, unit).After(time.Now()) {
		return false
	}
	return true
}

func calculateTargetDate(date time.Time, value int, unit string) time.Time {
	if uint == "d" {
		return date.AddDate(0, 0, value)
	}
	if uint == "w" {
		return date.AddDate(0, 0, value*7)
	}
	if unit == "m" {
		return date.AddDate(0, value, 0)
	}
	if unit == "y" {
		return date.AddDate(value, 0, 0)
	}
	return date
}
func hasCorrectPrefix(index string, prefixes []string) bool {
	result := false
	for _, prefix := range prefixes {
		if strings.HasPrefix(index, prefix) {
			return true
		}
	}
	return result
}

//索引名是否满足这个前缀
func trimPrefix(index string, prefixes []string) string {
	for _, prefix := range prefixes {
		if strings.HasPrefix(index, prefix) {
			return strings.TrimPrefix(index, prefix)
		}
	}
	return index
}
