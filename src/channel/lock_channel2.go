package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int)
	// val := <-c //it not right， locking
	go func() {
		time.Sleep(time.Second * 3)
		c <- 60
	}()
	val := <-c
	fmt.Println(val)
}

/* 总结

对同一个channel 的 read 和write 要在不同的goroutine中进行，在同一个goroutine中会死锁
从channel中的读取操作(<-c)代码顺序要在要在写入channel的操作之后(c<-) 要不也会产生死锁
从channel中读取操作会等待写入channel数据操作的执行，有数据了才能读嘛，channel读取操作会一直阻塞
*/
