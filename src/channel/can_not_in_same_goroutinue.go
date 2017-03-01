package main

import "fmt"

func main() {

	c := make(chan int)
	c <- 1
	a := <-c
	// c <- 1

	fmt.Println(a)
}

// By default, sends and receives block until the other side is ready
// so they can't be the same goroutinue, deadlock happend
