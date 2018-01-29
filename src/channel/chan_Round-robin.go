package main

import (
	"fmt"
	"runtime"
)

func main() {
	cores := runtime.NumCPU()
	runtime.GOMAXPROCS(cores)
	clients := make(chan int, 16)

	for i := 0; i < cores; i++ {
		clients <- i
	}

	for i := 0; i < 1000; i++ {
		client := <-clients
		clients <- client
		fmt.Println(client)
	}

}

// 输出的结果是 0 1 2 3 0 1 2 3 ..... 实现了简单的轮询效果
