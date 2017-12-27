package main

import "testing"

func BenchmarkChanCounter(b *testing.B) {
	c := chanCounter()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = <-c
	}
}

func BenchmarkMutexCounter(b *testing.B) {
	f := mutexCounter()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = f()
	}
}

func BenchmarkAtomicCounter(b *testing.B) {
	f := atomicCounter()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = f()
	}
}

// go test -v -bench . -benchmem
/*
goos: linux
goarch: amd64
BenchmarkChanCounter-4     	 5000000	       310 ns/op	       0 B/op	       0 allocs/op
BenchmarkMutexCounter-4    	50000000	        21.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkAtomicCounter-4   	200000000	         9.25 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	_/home/user/code/learning-go/src/optimize-go/channel	5.678s

*/
