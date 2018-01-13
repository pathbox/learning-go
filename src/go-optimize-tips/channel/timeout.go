package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	test()

	for i := 0; i < 10; i++ {
		runtime.GC()
		time.Sleep(time.Second)
	}
}

func test() {
	c := make(chan int) // 在goroutine 外部make chan

	go func() {
		for x := range c {
			fmt.Println(x)
		}
	}()

	for i := 0; i < 10; i++ {
		go func() { // select + timeout 超时机制,避免goroutine Leak
			select {
			case c <- 1:
			case <-time.After(time.Second * 2):
				return
			}
		}()
	}
}
