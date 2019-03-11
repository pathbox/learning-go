https://github.com/Xeoncross/go-cache-benchmark

It is from ``Xeoncross/go-cache-benchmark`

go get
go test -bench=. -benchmem


* https://golang.org/pkg/sync/#Map
* https://github.com/coocood/freecache
* https://github.com/allegro/bigcache
* https://github.com/patrickmn/go-cache
* https://github.com/muesli/cache2go
* https://github.com/bluele/gcache

go-cache-benchmark git:(master) ✗ go test -bench=. -benchmem
goos: darwin
goarch: amd64
BenchmarkCache2Go/Set-4                  1000000              2122 ns/op             370 B/op          9 allocs/op
BenchmarkCache2Go/Get-4                  3000000               582 ns/op              39 B/op          3 allocs/op
BenchmarkGoCache/Set-4                   1000000              1775 ns/op             279 B/op          5 allocs/op
BenchmarkGoCache/Get-4                   3000000               417 ns/op              23 B/op          2 allocs/op
BenchmarkFreecache/Set-4                 2000000              1053 ns/op              81 B/op          5 allocs/op
BenchmarkFreecache/Get-4                 5000000               441 ns/op              40 B/op          3 allocs/op
BenchmarkBigCache/Set-4                  1000000              1228 ns/op              99 B/op          4 allocs/op
BenchmarkBigCache/Get-4                  2000000               508 ns/op              47 B/op          3 allocs/op
BenchmarkGCache/Set-4                    1000000              2168 ns/op             207 B/op          8 allocs/op
BenchmarkGCache/Get-4                    5000000               285 ns/op              39 B/op          3 allocs/op
BenchmarkSyncMap/Set-4                   1000000              2396 ns/op             258 B/op          9 allocs/op
BenchmarkSyncMap/Get-4                   3000000               345 ns/op              23 B/op          2 allocs/op