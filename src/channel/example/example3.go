package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int)
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("before received")
		fmt.Println(<-c)
	}()

	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("before send")
		c <- 1
	}()

	time.Sleep(2 * time.Second)
	fmt.Println("after received")
}

// before received
// before send
// 1
// after received

// send is in the sub goroutinue
// send is be ready first, then received
// receive is not blocking. This is a nice way.
// but main goroutinue don't wait for two of them,
// c <- 1 block <-c, send first then received
// <-c wait for c<-1
