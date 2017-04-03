package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int)
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("before send")
		c <- 100
	}()

	fmt.Println(<-c)

	fmt.Println("after received")
}

// before send
// 100
// after received

// c <- in sub goroutinue, <- c in main goroutinue, c <- block main goroutinue
// It is the same as example1.go
// but no deadlock
