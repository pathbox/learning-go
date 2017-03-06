package main

import (
	"fmt"
)

func main() {
	naturals := make(chan int)
	squares := make(chan int)

	go func() {
		for x := 0; x < 100; x++ {
			naturals <- x
		}
		close(naturals)
	}()

	go func() {
		for x := range naturals {
			squares <- x * x
		}
		close(squares)
	}()

	for x := range squares {
		fmt.Println(x)
	}
}

// 发送完成后，可以关闭channel，关闭后所有对这个channel的写操作都会panic,而读操作依旧可以进行，当所有值都读完后，继续读该channel会得到zero value
