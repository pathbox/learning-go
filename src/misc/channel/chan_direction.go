package main

import (
	"fmt"
	"time"
)

type Sender chan<- int // send-only

type Receiver <-chan int //receive-only

func main() {
	var myChannel = make(chan int, 0)
	var number = 6
	go func() {
		var sender Sender = myChannel
		sender <- number
		fmt.Println("Sent!")
	}()

	go func() {
		var receiver Receiver = myChannel
		time.Sleep(time.Second)
		close(receiver)
		fmt.Println("Received!", <-receiver)
	}()

	// 让main函数执行结束的时间延迟1秒，
	// 以使上面两个代码块有机会被执行。
	time.Sleep(2 * time.Second)
}

// Received! 6
// Sent!

// number 进入 sender之后就阻塞了 直到<-receiver 的执行，然后chan就关闭了
// 所以 "Received!" 会在 "Sent!" 之前执行。如果第二个 func()不存在，sender会一直阻塞直到 main 的
// goroutine 执行结束而把整个进程关闭
