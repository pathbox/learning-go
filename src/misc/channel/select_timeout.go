package main

import (
	"fmt"
	"time"
)

func main() {
	w := make(chan bool)
	c := make(chan int, 2)

	go func() {
		select {
		case v := <-c:
			fmt.Println(v)
		case <-time.After(3* time.Second):
			fmt.Println("timeout.")
		}

		w <- true
	}()

	// c <- 1 // 注释掉，引发 timeout。
	<-w　// 用于阻塞外部goroutine　使得子goroutine go func()得以先完成
}　
