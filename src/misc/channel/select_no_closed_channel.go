package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	finish := make(chan bool)
	var done sync.WaitGroup
	done.Add(1)
	go func() {
		select { //select is a event-base loop until case condition is ok, it stop
		case <-time.After(1 * time.Hour):
		case <-finish:
		}
		done.Done()
	}()
	t0 := time.Now()
	finish <- true // 发送关闭信号
	done.Wait()    // 等待 goroutine 结束
	fmt.Printf("Waited %v for goroutine to stop\n", time.Since(t0))
}
