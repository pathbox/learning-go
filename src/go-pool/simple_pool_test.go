package simpool

import (
	"sync"
	"testing"
)

func Gopool() {
	wg := new(sync.WaitGroup)
	data := make(chan int, 100)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			for _ = range data {

			}
		}(i)
	}

	for i := 0; i < 10000; i++ {
		data <- i
	}
	close(data)
	wg.Wait()
}

func Nopool() {
	wg := new(sync.WaitGroup)

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
		}(i)
	}
	wg.Wait()
}

func BenchmarkGopool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Gopool()
	}
}

func BenchmarkNopool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Nopool()
	}
}

/* go test -v -bench . -benchmem
goos: linux
goarch: amd64
BenchmarkGopool-4           1000           1227901 ns/op             965 B/op          2 allocs/op
BenchmarkNopool-4            500           2571762 ns/op            4000 B/op         10 allocs/op
PASS
ok      _/home/user/code/learning-go/src/go-pool        2.924s
*/
