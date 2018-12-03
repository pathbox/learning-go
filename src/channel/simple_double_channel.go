package main

import (
	"fmt"
	"time"
)

func main() {
	cc := make(chan chan int)
	times := 5
	for i := 0; i < times; i++ {
		f := make(chan bool)

		// 每次循环都在双层通道cc中生成内层通道c
		// 并通过信号通道f来终止f1()
		go f1(cc, f)
		// 从双层通道cc中取出内层通道ch
		// 并向ch通道发送数据
		ch := <-cc // 得到内层通道
		ch <- i

		// 从ch中取出数据
		for sum := range ch {
			fmt.Printf("Sum(%d)=%d\n", i, sum)
		}
		time.Sleep(time.Second)
		// 每次循环都关闭信号通道f
		close(f)
	}
}

// 双层通道cc用来生成内层通道c
// 并使用信号通道f来终止函数f1()
func f1(cc chan chan int, f chan bool) {
	c := make(chan int) // new一个chan
	cc <- c             // 传入内层chan
	defer close(c)

	sum := 0
	select {
	// 从内层通道中取出数据，计算和，然后发回内层通道
	case x := <-c:
		for i := 0; i <= x; i++ {
			sum = sum + i
		}
		// goroutine将阻塞在此，直到数据被读走
		c <- sum
		// 信号通道f可读时，结束f1()的运行
		// 但因为select没有在for中，该case分支用不上
	case <-f:
		return
	}
}
