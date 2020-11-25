package concurrencyslower

import (
	"runtime"
	"sync"
	"testing"
)

const (
	limit = 10000000000
)

func ConcurrentSum() int {
	n := runtime.GOMAXPROCS(0)

	sums := make([]int, n)

	wg := sync.WaitGroup{}
	for i := 0; i < n; i++ {

		wg.Add(1)
		go func(i int) {
			start := (limit / n) * i
			end := start + (limit / n)

			for j := start; j < end; j += 1 {
				sums[i] += j
			}

			wg.Done()
		}(i)
	}

	wg.Wait()

	sum := 0
	for _, s := range sums {
		sum += s
	}
	return sum
}

func BenchmarkConcurrentSum(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ConcurrentSum()
	}
}

/*
go test -bench=. benchmark_test.go
go test -bench=. -cpuprofile=cpu.prof
go test -bench=. -memprofile=mem.prof
go tool pprof -http=:8080 cpu.prof
*/

/*
goos: darwin
goarch: amd64
BenchmarkSerialSum-4                   1        5573249889 ns/op
BenchmarkConcurrentSum-4               1        13457904182 ns/op
BenchmarkChannelSum-4                  1        2208850090 ns/op
PASS

https://appliedgo.net/concurrencyslower/

n具有n缓存的CPU内核重复读取和写入全部位于同一缓存行中的切片元素。因此，每当一个CPU内核用新的总和更新其“ slice”元素时，所有其他CPU的缓存行都会失效。更改后的高速缓存行必须写回到主存储器，所有其他高速缓存必须用新数据更新其各自的高速缓存行。即使每个核心都访问片的不同部分！

这会浪费宝贵的时间，而不是串行循环更新其单个sum变量所需的时间。

这就是为什么我们的并发循环比串行循环需要更多时间的原因。切片的所有并发更新都会导致高速缓存行同步跳动

既然我们知道了令人惊奇的减速的原因，那么治愈是显而易见的。我们必须将切片变成n单独的变量，希望它们彼此之间存储得足够远，以使它们不共享同一条缓存行。

因此，让我们更改并发循环，以便每个goroutine将其中间和存储在goroutine-local变量中。为了将结果传递回主goroutine，我们还必须添加一个通道。反过来，这使我们可以删除等待组，因为通道不仅是一种通信手段，而且还是一种优雅的同步机制。
*/
