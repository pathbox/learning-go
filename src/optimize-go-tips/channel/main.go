package main

import (
	"sync"
	"sync/atomic"
)

func chanCounter() chan int {
	c := make(chan int)
	go func() {
		for x := 1; ; x++ {
			c <- x
		}
	}()
	return c
}

func mutexCounter() func() int {
	var m sync.Mutex
	var x int

	return func() (n int) {
		m.Lock()
		x++
		n = x
		m.Unlock()
		return
	}
}

func atomicCounter() func() int {
	var x int64

	return func() int {
		return int(atomic.AddInt64(&x, 1))
	}
}

func main() {
	c := chanCounter()
	println(<-c)

	m := mutexCounter()
	println(m())
}
