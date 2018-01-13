package main

import "testing"

func BenchmarkTest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		data := make(chan int, bufsize)
		done := make(chan struct{})
		b.StartTimer()

		_ = test(data, done)
	}
}

func BenchmarkBlockTest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		data := make(chan [block]int, bufsize)
		done := make(chan struct{})
		b.StartTimer()

		_ = testBlock(data, done)
	}
}

func BenchmarkSliceBlockTest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		data := make(chan []int, bufsize)
		done := make(chan struct{})
		b.StartTimer()

		_ = testSliceBlock(data, done)
	}
}

/*
goos: linux
goarch: amd64
BenchmarkTest-4             	      30	  41778975 ns/op	      35 B/op	       1 allocs/op
BenchmarkBlockTest-4        	    1000	   2494043 ns/op	      32 B/op	       1 allocs/op
BenchmarkSliceBlockTest-4   	     500	   2170738 ns/op	 4096045 B/op	    1001 allocs/op
PASS
ok  	_/home/user/code/learning-go/src/optimize-go-tips/channel/more	5.712s

slice 非但没有提升性能，反而在堆上分配了更多内存，有些得不偿失
*/
