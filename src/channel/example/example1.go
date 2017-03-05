package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int)
	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println("before received")
		// fmt.Println(<-c) // 这里在阻塞， 这里会先执行 ready
	}()

	c <- 1
	fmt.Println("after received")
}

// send in the main goroutinue
// before received
// 1
// after received
// main is blocked，until channel is read，then main go on  因为写进channel 的数据需要被读取出来
// 如果读取操作在子goroutinue， 写操作在main goroutinue， 读取操作会需要先执行，并且阻塞main goroutinue
// 如果 channel 读取操作和写操作 都在子goroutinue 则谁先就谁先ready
// 见example4.go  和 example5.go
