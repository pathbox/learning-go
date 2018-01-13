package main

import "testing"

func BenchmarkUseMutex(b *testing.B) {
	for n := 0; n < b.N; n++ {
		UseMutex()
	}
}
func BenchmarkUseChan(b *testing.B) {
	for n := 0; n < b.N; n++ {
		UseChan()
	}
}

/* go test -bench=.
goos: linux
goarch: amd64
BenchmarkUseMutex-4     100000000               19.0 ns/op
BenchmarkUseChan-4      20000000                60.7 ns/op
PASS
*/
