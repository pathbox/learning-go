package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int)
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("before received")
		fmt.Println("print channel", <-c) // 这里在阻塞， 这里会先执行 ready
	}()

	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println("before send")
		c <- 1
	}()

	time.Sleep(1 * time.Second)

	fmt.Println("after received")
}

// after received
// main goroutinue will not wait for two sub goroutinue, they will be killed when main goroutinue is over
