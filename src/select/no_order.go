// select  的case是无序的，当有多个case同时满足时，调度器选择执行任意一个都有可能，
// 所以，要在代码层面保证select 的case 在同一时刻只有一个满足，除非觉得任意个执行都可以
package main

import (
	"fmt"
	"time"
)

func main() {
	c1 := make(chan int, 10)
	c2 := make(chan int, 10)

	go func() {
		for i := 0; i < 100; i++ {
			c1 <- 1
		}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			c2 <- 1
		}
	}()

	time.Sleep(100 * time.Millisecond)

	for {
		select {
		case i := <-c1:
			fmt.Println("c1: ", i)
		case i := <-c2:
			fmt.Println("c2: ", i)
		}
	}
}

/*
c1:  1
c1:  1
c1:  1
c2:  1
c2:  1
c1:  1
c2:  1
c1:  1
c2:  1
c1:  1
c1:  1
c1:  1
c1:  1
c1:  1
c1:  1
c2:  1
c2:  1
c2:  1
c2:  1
c2:  1
c2:  1
c2:  1
c2:  1
...
*/
