package main

import (
	"fmt"
)

func main() {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	// 会报死锁错误的，原因是range不等到信道关闭是不会结束读取的。也就是如果 缓冲信道干涸了，那么range就会阻塞当前goroutine, 所以死锁咯

	for v := range ch {
		fmt.Println(v)
		if len(ch) <= 0 {
			close(ch)
		}
	}
}
