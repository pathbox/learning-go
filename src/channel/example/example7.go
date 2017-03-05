package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int)
	go func() {

		fmt.Println("before received")
		fmt.Println(<-c) // 这里在阻塞，直到 数据过来
	}()
	time.Sleep(1 * time.Second)
	fmt.Println("before send")
	c <- 1
	time.Sleep(2 * time.Second)
	fmt.Println("after received")
}

// fmt.Println(<-c) <-c is blocking, until c <- 1
// 新开的goroutine首先去读channel,可是由于channel中没有值，所以它被阻塞了，直到main中向channel发送值，goroutine才拿到它想要的值并继续运行。
