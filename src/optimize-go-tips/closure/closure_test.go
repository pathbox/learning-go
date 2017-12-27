package main

import (
	"testing"
)

func test(x int) int {
	return x * 2
}

func BenchmarkTest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = test(i)
	}
}

func BenchmarkClosure(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = func() int {
			return i * 2
		}()
	}
}

func BenchmarkAnonymous(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = func(x int) int {
			return x * 2
		}(i)
	}
}

// go test -v -bench . -benchmem
/*
goos: linux
goarch: amd64
BenchmarkTest-4        	2000000000	         0.33 ns/op	       0 B/op	       0 allocs/op
BenchmarkClosure-4     	1000000000	         2.60 ns/op	       0 B/op	       0 allocs/op
BenchmarkAnonymous-4   	1000000000	         2.27 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	_/home/user/code/learning-go/src/optimize-go/closure	6.095s

*/

// https://www.jianshu.com/p/eb82216ea8ab
