package main

import "testing"

func BenchmarkIncX(b *testing.B) {
	d := struct {
		X int
	}{100}

	for i := 0; i < b.N; i++ {
		incX(&d)
	}
}

func BenchmarkUnsafeIncX(b *testing.B) {
	d := struct {
		X int
	}{100}

	for i := 0; i < b.N; i++ {
		unsafeIncX(&d)
	}
}

func BenchmarkCacheUnsafeIncX(b *testing.B) {
	d := struct {
		X int
	}{100}

	for i := 0; i < b.N; i++ {
		unsafeCacheIncX(&d)
	}
}

/*
goos: linux
goarch: amd64
BenchmarkIncX-4              	10000000	       139 ns/op	       8 B/op	       1 allocs/op
BenchmarkUnsafeIncX-4        	500000000	         3.86 ns/op	       0 B/op	       0 allocs/op
BenchmarkCacheUnsafeIncX-4   	100000000	        17.0 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	_/home/user/code/learning-go/src/optimize-go/reflect	5.591s

利用指针类型转换实现性能优化，本就是 “非常手段”，是一种为了性能而放弃 “其他” 的做法。
与其担心代码是否适应未来的变化，不如写个单元测试，确保在升级时做出必要的安全检查
*/

// https://www.jianshu.com/p/d86529279e1e
