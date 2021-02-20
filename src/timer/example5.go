package main

import (
	"time"
	"fmt"
)

func main() {
	timer := time.NewTimer(time.Second)
	c := make(chan bool)

	go func() {
		for {
			time.Sleep(time.Second * 2)
			c <- false
		}
	}()

	for {
		select {
		case <- timer.C:
			fmt.Println("timer")
			timer.Reset(time.Second) // 这里使用NewTimer定时器需要t.Reset重置计数时间才能接着执行,要不只会执行一次
		case val := <-c:
			fmt.Println(val)
		}
	}
	fmt.Println("done")
}