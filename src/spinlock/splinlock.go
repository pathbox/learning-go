package spinLock

import (
	"runtime"
	"sync/atomic"
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
