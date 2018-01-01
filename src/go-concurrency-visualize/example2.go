package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 0; i < 24; i++ {
		c := tick(100 * time.Millisecond)
		fmt.Println(<-c) // <- chan
	}
}

func tick(d time.Duration) <-chan int {
	c := make(chan int)
	go func() {
		time.Sleep(d)
		c <- 1
	}()
	return c
}
