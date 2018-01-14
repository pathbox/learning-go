package main

import (
	"testing"
)

func test(m map[int]int) {
	for i := 0; i < 10000; i++ {
		m[i] = i
	}
}

func BenchmarkMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		m := make(map[int]int)
		b.StartTimer()

		test(m)
	}
}

func BenchmarkCapMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		m := make(map[int]int, 10000)
		b.StartTimer()

		test(m)
	}
}

// 预设容量的 map 显然性能更好，更极大减少了堆内存分配次数

/*
go test -v -bench . -benchmem
goos: linux
goarch: amd64
BenchmarkMap-4      	    1000	   1446738 ns/op	  687223 B/op	     276 allocs/op
BenchmarkCapMap-4   	    3000	    619656 ns/op	    2721 B/op	       9 allocs/op
PASS
ok  	_/home/user/code/learning-go/src/optimize-go/map	3.808s

*/
