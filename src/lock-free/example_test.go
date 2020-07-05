package example_test

import (
	"sync"
	"sync/atomic"
	"testing"
)

type Config struct {
	sync.RWMutex
	endpoint string
}

func BenchmarkPMutexSet(b *testing.B) {
	config := Config{}
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			config.Lock()
			config.endpoint = "api.example.com"
			config.Unlock()
		}
	})
}

func BenchmarkPMutexGet(b *testing.B) {
	config := Config{endpoint: "api.example.com"}
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			config.RLock()
			_ = config.endpoint
			config.RUnlock()
		}
	})
}

func BenchmarkPAtomicSet(b *testing.B) {
	var config atomic.Value
	c := Config{endpoint: "api.example.com"}
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			config.Store(c)
		}
	})
}

func BenchmarkPAtomicGet(b *testing.B) {
	var config atomic.Value
	config.Store(Config{endpoint: "api.example.com"})
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = config.Load().(Config)
		}
	})
}

/*
go test -v -bench=. -benchmem example_test.go
goos: darwin
goarch: amd64
BenchmarkPMutexSet
BenchmarkPMutexSet-4            11223747               106 ns/op               0 B/op          0 allocs/op
BenchmarkPMutexGet
BenchmarkPMutexGet-4            31274442                49.5 ns/op             0 B/op          0 allocs/op
BenchmarkPAtomicSet
BenchmarkPAtomicSet-4           20757464                62.9 ns/op            48 B/op          1 allocs/op
BenchmarkPAtomicGet
BenchmarkPAtomicGet-4           1000000000               1.01 ns/op            0 B/op          0 allocs/op
PASS
 Get 的性能相比惊人
*/
