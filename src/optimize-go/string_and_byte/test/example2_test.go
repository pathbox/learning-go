package main

import (
	"strings"
	"testing"
)

var s = strings.Repeat("a", 1024)

// 字符串转为byte,byte再转为字符串
func test() {
	b := []byte(s)
	_ = string(b)
}

func test2() {
	b := string2bytes(s)
	_ = bytes2string(b)
}

func BenchmarkTest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		test()
	}
}

func BenchmarkTestNice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		test2()
	}
}

/*
go test -v -bench . -benchmem

goos: linux
goarch: amd64
BenchmarkTest-4        	 3000000	       484 ns/op	    2048 B/op	       2 allocs/op
BenchmarkTestNice-4   	500000000	         3.85 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	_/home/user/code/learning-go/src/optimize-go/string_and_byte/test	4.234s
*/

/* string转为byte的benchmark 性能有巨大差别 性能提升明显，最关键的是 zero-garbage
goos: linux
goarch: amd64
BenchmarkTest-4        	 5000000	       240 ns/op	    1024 B/op	       1 allocs/op
BenchmarkTestNice-4   	1000000000	         2.63 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	_/home/user/code/learning-go/src/optimize-go/string_and_byte/test	4.363s
 ~/code/learning-go/src/optimize-go/string_and_byte/test   master  go test -v -bench . -benchmem
byte转为string的benchmark 性能没什么问题
 goos: linux
goarch: amd64
BenchmarkTest-4        	100000000	        13.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkTestNice-4   	200000000	         9.02 ns/op	       0 B/op	       0 allocs/op
PASS
*/
