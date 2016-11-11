package main

import (
	"fmt"
)

// 缓冲信道是先进先出的，我们可以把缓冲信道看作为一个线程安全的队列

func main() {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3

	fmt.Println(<-ch) // 1
	fmt.Println(<-ch) // 2
	fmt.Println(<-ch) // 3
}
