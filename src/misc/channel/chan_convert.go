package main

import (
	"fmt"
)

func main() {
	c := make(chan int, 3)

	// 可以将 channel 隐式转换为单向队列，只收或只发
	var send chan<- int = c // send-only
	var recv <-chan int = c // receive-only
	send <- 1
	// <-send // Error: receive from send-only type chan<- int
	fmt.Println(<-recv)
	// recv <- 2 // Error: send to receive-only type <-chan int
	close(c)
	// 不能将单向 channel 转换为普通 channel
	// d := (chan int)(send) // Error: cannot convert type chan<- int to type chan int
	// d := (chan int)(recv) // Error: cannot convert type <-chan int to type chan int
}
