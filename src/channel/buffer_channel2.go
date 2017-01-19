package main

import (
	"fmt"
	"time"
)

func main() {
	message := make(chan string, 2)
	count := 3

	go func() {
		for i := 1; i <= count; i++ {
			fmt.Println("send message")
			message <- fmt.Sprintf("message %d", i)
		}
	}()

	time.Sleep(time.Second * 3)
	for i := 1; i <= count; i++ {
		fmt.Println(<-message)
	}
}
