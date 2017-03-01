package main

import (
	"fmt"
)

var c = make(chan int)

// var c = make(chan int, 1) // 最后可能打印的是空

// channel c guarante a = "hello world" before fmt.Println(a)
var a string

func f() {
	a = "hello world"
	<-c
}

func main() {
	go f()
	c <- 0
	fmt.Println(a)
}
