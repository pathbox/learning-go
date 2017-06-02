// This is a good way to set a sync.Mutex in the struct. It is the struct Mutext to lock the struct's data
package main

import (
	"fmt"
	"sync"
)

type safeCounter struct {
	number int
	sync.Mutex
}

func (sc *safeCounter) Increment() {
	sc.Lock()
	sc.number++
	sc.Unlock()
}

func (sc *safeCounter) Decrement() {
	sc.Lock()
	sc.number--
	sc.Unlock()
}

func (sc *safeCounter) getNumber() int {
	sc.Lock()
	number = sc.number
	sc.Unlock()
	return number
}

func main() {
	sc := new(safeCounter)
	for i := 0; i < 100; i++ {
		go sc.Increment()
		go sc.Decrement()
	}
	fmt.Println(sc.getNumber())
}
