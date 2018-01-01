package main

import (
	"fmt"
)

func main() {
	ch := make(chan int)
	go func() {
		ch <- 42
	}()

	fmt.Println(<-ch)
}
