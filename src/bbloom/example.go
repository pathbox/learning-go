package main

func main() {
	// create a bloom filter for 65536 items and 1 % wrong-positive ratio
	bf := bbloom.New(float64(1<<16), float64(0.01))

	// or
	// create a bloom filter with 650000 for 65536 items and 7 locs per hash explicitly
	// bf = bbloom.New(float64(650000), float64(7))
	// or
	bf = bbloom.New(650000.0, 7.0)

	// add one item
	bf.Add([]byte("butter"))

	// Number of elements added is exposed now
	// Note: ElemNum will not be included in JSON export (for compatability to older version)
	nOfElementsInFilter := bf.ElemNum

	// check if item is in the filter
	isIn := bf.Has([]byte("butter"))    // should be true
	isNotIn := bf.Has([]byte("Butter")) // should be false

	// 'add only if item is new' to the bloomfilter
	added := bf.AddIfNotHas([]byte("butter")) // should be false because 'butter' is already in the set
	added = bf.AddIfNotHas([]byte("buTTer"))  // should be true because 'buTTer' is new

	// thread safe versions for concurrent use: AddTS, HasTS, AddIfNotHasTS
	// add one item
	bf.AddTS([]byte("peanutbutter"))
	// check if item is in the filter
	isIn = bf.HasTS([]byte("peanutbutter"))    // should be true
	isNotIn = bf.HasTS([]byte("peanutButter")) // should be false
	// 'add only if item is new' to the bloomfilter
	added = bf.AddIfNotHasTS([]byte("butter"))       // should be false because 'peanutbutter' is already in the set
	added = bf.AddIfNotHasTS([]byte("peanutbuTTer")) // should be true because 'penutbuTTer' is new

	// convert to JSON ([]byte)
	Json := bf.JSONMarshal()

	// bloomfilters Mutex is exposed for external un-/locking
	// i.e. mutex lock while doing JSON conversion
	bf.Mtx.Lock()
	Json = bf.JSONMarshal()
	bf.Mtx.Unlock()

	// restore a bloom filter from storage
	bfNew := bbloom.JSONUnmarshal(Json)

	isInNew := bfNew.Has([]byte("butter"))    // should be true
	isNotInNew := bfNew.Has([]byte("Butter")) // should be false
}
