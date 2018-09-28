package main

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

const (
	unlocked int32 = iota //0
	locked                // 1
)

// A spinLock must not be copied after first use.
type spinLock struct {
	state int32 // 初始是0 unlocked
}

func (lock *spinLock) Lock() {
	for !atomic.CompareAndSwapInt32(&lock.state, unlocked, locked) {
		runtime.Gosched()
	}
}

func (lock *spinLock) Unlock() {
	for !atomic.CompareAndSwapInt32(&lock.state, locked, unlocked) {
		runtime.Gosched()
	}
}
func BenchmarkLock(b *testing.B) {
	var lock spinLock

	for i := 0; i < b.N; i++ {
		lock.Lock()
		lock.Unlock()
	}
}

func BenchmarkMutex(b *testing.B) {
	var lock sync.Mutex

	for i := 0; i < b.N; i++ {
		lock.Lock()
		lock.Unlock()
	}
}

func BenchmarkWMutex(b *testing.B) {
	var lock sync.RWMutex

	for i := 0; i < b.N; i++ {
		lock.Lock()
		lock.Unlock()
	}
}

func BenchmarkRMutex(b *testing.B) {
	var lock sync.RWMutex

	for i := 0; i < b.N; i++ {
		lock.RLock()
		lock.RUnlock()
	}
}

/*
go test -bench=. -cpu 1,2,3 locks_test.go
goos: darwin
goarch: amd64
BenchmarkLock           100000000               16.0 ns/op
BenchmarkLock-2         100000000               16.0 ns/op
BenchmarkLock-3         100000000               16.0 ns/op
BenchmarkMutex          100000000               16.0 ns/op
BenchmarkMutex-2        100000000               16.0 ns/op
BenchmarkMutex-3        100000000               16.1 ns/op
BenchmarkWMutex         50000000                38.3 ns/op
BenchmarkWMutex-2       50000000                34.8 ns/op
BenchmarkWMutex-3       50000000                34.3 ns/op
BenchmarkRMutex         100000000               16.1 ns/op
BenchmarkRMutex-2       100000000               16.3 ns/op
BenchmarkRMutex-3       100000000               16.1 ns/op
PASS
*/
