package main

import (
	"runtime"
	"sync/atomic"
	"testing"
	"time"
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

func TestLock(t *testing.T) {
	mx := &spinLock{}
	resource := make(map[int]int)
	done := make(chan struct{})

	go func() {
		for i := 0; i < 10; i++ {
			mx.Lock()
			resource[i] = i
			time.Sleep(time.Millisecond)
			mx.Unlock()
		}

		done <- struct{}{}
	}()

	for i := 0; i < 10; i++ {
		mx.Lock()
		_ = resource[i]
		time.Sleep(time.Millisecond)
		mx.Unlock()
	}

	<-done
}
