package ants_test

import (
	"sync"
	"testing"

	"github.com/panjf2000/ants"
)

const (
	RunTimes      = 1000000
	benchParam    = 10
	benchAntsSize = 100000
	m             = 100000
)

func demoFunc() error {
	// n := 10
	// time.Sleep(time.Duration(n) * time.Millisecond)
	var n int
	for i := 0; i < m; i++ {
		n += i
	}
	return nil
}

func demoPoolFunc(args interface{}) error {
	var n int
	for i := 0; i < m; i++ {
		n += i
	}
	// n := args.(int)
	// time.Sleep(time.Duration(n) * time.Millisecond)
	return nil
}

func BenchmarkGoroutineWithFunc(b *testing.B) {
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(RunTimes)
		for j := 0; j < RunTimes; j++ {
			go func() {
				demoPoolFunc(benchParam)
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkSemaphoreWithFunc(b *testing.B) {
	var wg sync.WaitGroup
	sema := make(chan struct{}, benchAntsSize)

	for i := 0; i < b.N; i++ {
		wg.Add(RunTimes)
		for j := 0; j < RunTimes; j++ {
			sema <- struct{}{}
			go func() {
				demoPoolFunc(benchParam)
				<-sema
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkAntsPoolWithFunc(b *testing.B) {
	var wg sync.WaitGroup
	p, _ := ants.NewPoolWithFunc(benchAntsSize, func(i interface{}) error {
		demoPoolFunc(i)
		wg.Done()
		return nil
	})
	defer p.Release()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(RunTimes)
		for j := 0; j < RunTimes; j++ {
			p.Serve(benchParam)
		}
		wg.Wait()
		//b.Logf("running goroutines: %d", p.Running())
	}
	b.StopTimer()
}

func BenchmarkGoroutine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < RunTimes; j++ {
			go demoPoolFunc(benchParam)
		}
	}
}

func BenchmarkSemaphore(b *testing.B) {
	sema := make(chan struct{}, benchAntsSize)
	for i := 0; i < b.N; i++ {
		for j := 0; j < RunTimes; j++ {
			sema <- struct{}{}
			go func() {
				demoPoolFunc(benchParam)
				<-sema
			}()
		}
	}
}

func BenchmarkAntsPool(b *testing.B) {
	p, _ := ants.NewPoolWithFunc(benchAntsSize, demoPoolFunc)
	defer p.Release()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < RunTimes; j++ {
			p.Serve(benchParam)
		}
	}
	b.StopTimer()
}

/*
--- pool size: 10w
almost same

--- pool size: 100w goroutine size: 100w
go test -bench="GoroutineWithFunc" -benchmem=true -run=none
goos: darwin
goarch: amd64
pkg: github.com/pathbox/learning-go/src/pool/ants_pool/ants_test
BenchmarkGoroutineWithFunc-4   	       1	22866796601 ns/op	375096816 B/op	  882853 allocs/op
PASS
ok  	github.com/pathbox/learning-go/src/pool/ants_pool/ants_test	24.952s

go test -bench="AntsPoolWithFunc" -benchmem=true -run=none
goos: darwin
goarch: amd64
pkg: github.com/pathbox/learning-go/src/pool/ants_pool/ants_test
BenchmarkAntsPoolWithFunc-4   	       1	18684239454 ns/op	30559328 B/op	  215991 allocs/op
PASS
ok  	github.com/pathbox/learning-go/src/pool/ants_pool/ants_test	18.722s


go test -bench="Goroutine$" -benchmem=true -run=none
goos: darwin
goarch: amd64
pkg: github.com/pathbox/learning-go/src/pool/ants_pool/ants_test
BenchmarkGoroutine-4   	       1	28372360610 ns/op	374957312 B/op	  882488 allocs/op
PASS
ok  	github.com/pathbox/learning-go/src/pool/ants_pool/ants_test	32.684s

go test -bench="AntsPool$" -benchmem=true -run=none
goos: darwin
goarch: amd64
pkg: github.com/pathbox/learning-go/src/pool/ants_pool/ants_test
BenchmarkAntsPool-4   	       1	18968215077 ns/op	26878496 B/op	  187864 allocs/op
PASS
ok  	github.com/pathbox/learning-go/src/pool/ants_pool/ants_test	19.214s

pool size: 10w do time 100w
go test -bench="AntsPool$" -benchmem=true -run=none
goos: darwin
goarch: amd64
pkg: github.com/pathbox/learning-go/src/pool/ants_pool/ants_test
BenchmarkAntsPool-4   	       1	18707930244 ns/op	25872960 B/op	  180008 allocs/op
PASS
ok  	github.com/pathbox/learning-go/src/pool/ants_pool/ants_test	18.745s

go test -bench="AntsPoolWithFunc" -benchmem=true -run=none
goos: darwin
goarch: amd64
pkg: github.com/pathbox/learning-go/src/pool/ants_pool/ants_test
BenchmarkAntsPoolWithFunc-4   	       1	18751178572 ns/op	24489056 B/op	  169197 allocs/op
PASS
ok  	github.com/pathbox/learning-go/src/pool/ants_pool/ants_test	18.784s

*/
