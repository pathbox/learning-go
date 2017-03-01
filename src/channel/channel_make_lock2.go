package main

import (
	"fmt"
)

var c = make(chan int, 10)
var a string

func f() {
	a = "Hello world"
	c <- 0
}

func main() {
	go f()
	<-c // receive 会一直阻塞直到别的goroutine 进行了 c <- send操作才会继续往下执行，在进行send操作之前 a = "Hello world" 被执行了
	fmt.Println(a)
}

// 当满足下面条件时可以保证读操作r能侦测到写操作w：
// w happens-before r.
// Any other write to the shared variable v either happens-before w or after r.
// 关于channel的happens-before在Go的内存模型中提到了三种情况：
// 对一个channel的发送操作 happens-before 相应channel的接收操作完成
// 关闭一个channel happens-before 从该Channel接收到最后的返回值0
// 不带缓冲的channel的接收操作 happens-before 相应channel的发送操作完成
