package benchmark_test

import (
	// "fmt"
	"github.com/mohong122/ip2region/binding/golang"
	"testing"
)

func BenchmarkMemorySearch(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Search()
	}
}

func Search() {
	region, err := ip2region.New(".././ip2region.db")
	defer region.Close()
	if err != nil {
		return
	}
	// region.BinarySearch("123.95.223.18")
	// region.BtreeSearch("123.95.223.90")
	region.MemorySearch("123.95.223.90")
}

/*
go test -bench=.


 BtreeSearch

BenchmarkMemorySearch-4       100000       16815 ns/op
PASS
ok    _/home/user/Documents/udesk_vistor_go/benchmark 1.856s




MemorySearch

 2000     766670 ns/op
PASS
ok    _/home/user/Documents/udesk_vistor_go/benchmark 1.610s


BinarySearch

BenchmarkMemorySearch-4       100000       15430 ns/op
PASS
ok    _/home/user/Documents/udesk_vistor_go/benchmark 1.711s
*/
