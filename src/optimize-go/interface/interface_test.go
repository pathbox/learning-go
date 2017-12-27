package main

import "testing"

func BenchmarkCall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		call(&Data{x: 100})
	}
}

func BenchmarkIfaceCall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ifaceCall(&Data{x: 100})
	}
}

/*
go test -v -bench . -benchmem
goos: linux
goarch: amd64
BenchmarkCall-4        	2000000000	         0.32 ns/op	       0 B/op	       0 allocs/op
BenchmarkIfaceCall-4   	50000000	        31.0 ns/op	       8 B/op	       1 allocs/op
PASS
ok  	_/home/user/code/learning-go/src/optimize-go/interface	2.260s
*/

// 对于压力很大的内部组件之间，用接口有些得不偿失
// 普通调用被内联，但接口调用就没有这个待遇了
