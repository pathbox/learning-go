package main

import (
	"fmt"
	"time"
)

const (
	fmat = "2006-01-02 15:04:05"
)

func main() {
	c := make(chan int)
	go func() {
		time.Sleep(1 * time.Second)
		close(c)
	}()

	for {
		select {
		case x, ok := <-c:
			fmt.Printf("%v 通道读取到: x=%v ok=%v\n", time.Now().Format(fmat), x, ok)
			time.Sleep(500 * time.Millisecond)
		default:
			fmt.Printf("%v 通道没有读取到\n", time.Now().Format(fmat))
			time.Sleep(500 * time.Millisecond)
		}
	}
}
