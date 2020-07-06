https://tonybai.com/2020/04/24/gogoprotobuf-vs-goprotobuf-v1-and-v2/

```
make benchmark
cd goprotobuf && go test -bench .
goos: darwin
goarch: amd64
pkg: goprotobuf
BenchmarkMarshal-4             	 1535812	       768 ns/op	     384 B/op	       1 allocs/op
BenchmarkUnmarshal-4           	 1748592	       679 ns/op	     400 B/op	       7 allocs/op
BenchmarkMarshalInParalell-4   	 3401350	       355 ns/op	     384 B/op	       1 allocs/op
BenchmarkUnmarshalParalell-4   	 2762788	       437 ns/op	     400 B/op	       7 allocs/op
PASS
ok  	goprotobuf	7.082s
cd gogoprotobuf-fast && go test -bench .
goos: darwin
goarch: amd64
pkg: goprotobuffast
BenchmarkMarshal-4             	 4800471	       240 ns/op	     384 B/op	       1 allocs/op
BenchmarkUnmarshal-4           	 3310730	       354 ns/op	     400 B/op	       7 allocs/op
BenchmarkMarshalInParalell-4   	 7425297	       164 ns/op	     384 B/op	       1 allocs/op
BenchmarkUnmarshalParalell-4   	 4110501	       282 ns/op	     400 B/op	       7 allocs/op
PASS
ok  	goprotobuffast	5.796s
cd gogoprotobuf-faster && go test -bench .
goos: darwin
goarch: amd64
pkg: goprotobuffaster
BenchmarkMarshal-4             	 4716014	       241 ns/op	     384 B/op	       1 allocs/op
BenchmarkUnmarshal-4           	 3202999	       359 ns/op	     400 B/op	       7 allocs/op
BenchmarkMarshalInParalell-4   	 7455350	       158 ns/op	     384 B/op	       1 allocs/op
BenchmarkUnmarshalParalell-4   	 3601984	       297 ns/op	     400 B/op	       7 allocs/op
PASS
ok  	goprotobuffaster	5.693s
cd gogoprotobuf-slick && go test -bench .
goos: darwin
goarch: amd64
pkg: goprotobufslick
BenchmarkMarshal-4             	 4902349	       245 ns/op	     384 B/op	       1 allocs/op
BenchmarkUnmarshal-4           	 3369596	       374 ns/op	     400 B/op	       7 allocs/op
BenchmarkMarshalInParalell-4   	 7597167	       162 ns/op	     384 B/op	       1 allocs/op
BenchmarkUnmarshalParalell-4   	 4021544	       318 ns/op	     400 B/op	       7 allocs/op
PASS
ok  	goprotobufslick	6.055s

```