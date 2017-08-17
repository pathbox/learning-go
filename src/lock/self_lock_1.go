package main

import (
	"fmt"

	"time"
)

type record struct {
	lock *mutex
}

type mutex struct {
	lock chan struct{}
}

func newMutex() *mutex {
	return &mutex{lock: make(chan struct{}, 1)}
}

func newRecord() *record {
	return &record{
		lock: newMutex(),
	}
}

func (m *mutex) Lock() {
	m.lock <- struct{}{}
}

func (m *mutex) Unlock() {
	<-m.lock
}

func main() {
	r := newRecord()
	go hello_no_lock()
	go hello_lock(r)
	go hello_lock(r)
	go hello_no_lock()
	time.Sleep(10 * time.Second)
	fmt.Println("done")
}

func hello_lock(r *record) {
	r.lock.Lock()
	defer r.lock.Unlock()
	fmt.Println("hello lock")
	time.Sleep(5 * time.Second)

}

func hello_no_lock() {
	fmt.Println("hello no lock")
	time.Sleep(100 * time.Second)
}
