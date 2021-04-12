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
		c <- 10
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

/*
c通道是一个缓冲为0的通道，在main开始时，启动一个协程对c通道写入10，然后就关闭掉这个通道。
在main中通过 x, ok := <-c 接受通道c里的值，从输出结果里看出，确实从通道里读出了之前塞入通道的10，但是在通道关闭后，这个通道一直能读出内容。
*/
