// TryLock() if the lock is released， it returns false， if it hold the lock， it returns true

package trylock

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

const mutexLocked = 1 << iota

// Mutex is simple sync.Mutex + ability to try to Lock.
type Mutex struct {
	in sync.Mutex
}

// Lock locks m.
// If the lock is already in use, the calling goroutine
// blocks until the mutex is available.
func (m *Mutex) Lock() {
	m.in.Lock()
}

// Unlock unlocks m.
// It is a run-time error if m is not locked on entry to Unlock.
//
// A locked Mutex is not associated with a particular goroutine.
// It is allowed for one goroutine to lock a Mutex and then
// arrange for another goroutine to unlock it.
func (m *Mutex) Unlock() {
	m.in.Unlock()
}

// TryLock tries to lock m. It returns true in case of success, false otherwise.
func (m *Mutex) TryLock() bool {
	return atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.in)), 0, mutexLocked)
}
