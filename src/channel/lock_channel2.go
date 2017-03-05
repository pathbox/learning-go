package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int)
	// val := <-c //it not right， locking
	go func() {
		time.Sleep(time.Second * 1)
		c <- 60
	}()
	val := <-c
	fmt.Println(val)
}

/* 总结

对同一个channel 的 read 和write 要在不同的goroutine中进行，在同一个goroutine中会死锁
从channel中读取操作会等待写入channel数据操作的执行，有数据了才能读嘛，channel读取操作会一直阻塞
然后 <-c 操作阻塞等待数据的到来。
所以，建议尽快的让穿c<- send操作执行，这样 <-c 就不会阻塞
*/

/*
先进行写channel 的 goroutine代码，再进行读 channel的goroutine代码，读写逻辑分别在不同的goroutine；读channel会一直阻塞读的goroutine，直到读取出数据，除非用缓存channel。 读channel的goroutine知道有写channel的goroutine存在，然后就一直等数据过来，读取到之后，才结束阻塞
*/
