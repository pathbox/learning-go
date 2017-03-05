package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int)
	go func() {
		fmt.Println("before send")
		c <- 100
		fmt.Println("after send")

	}()

	time.Sleep(1 * time.Second)
	fmt.Println("done")
}

// before send
// done

// no receive, send c <- 100 in sub goroutinue
// fmt.Println("after send") don't run
// no deadlock
