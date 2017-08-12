package main

import (
	"./locker"
	"fmt"
	"time"
)

type record struct {
	lock          *locker.Mutex
	lock_count    int
	no_lock_count int
}

func newRecord() *record {
	return &record{
		lock:          locker.NewMutex(),
		lock_count:    0,
		no_lock_count: 0,
	}
}

func main() {
	r := newRecord()

	for i := 0; i < 1000; i++ {
		go CountWithoutLock(r)
		go CountWithLock(r)
	}
	time.Sleep(2 * time.Second)
	fmt.Println("Record no_lock_count: ", r.no_lock_count)
	fmt.Println("Record lock_count: ", r.lock_count)
}

func CountWithLock(r *record) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.lock_count++
}

func CountWithoutLock(r *record) {
	r.no_lock_count++
}

// 输出结果
// Record no_lock_count:  991
// Record lock_count:  1000
