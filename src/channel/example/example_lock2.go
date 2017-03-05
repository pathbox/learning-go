package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int)
	go func() {
		fmt.Println("before received")
		fmt.Println("receive", <-c)
		fmt.Println("after received")

	}()

	time.Sleep(1 * time.Second)
	fmt.Println("done")
}

// before received
// done

// no send, receive will get nil, receive <-c in sub goroutinue
// fmt.Println("receive", <-c) and fmt.Println("after received") don't run
// no deadlock
