package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	quit := make(chan bool)

	for i := 0; i <= 2; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			task := func() {
				fmt.Println(id, time.Now().Nanosecond())
				time.Sleep(time.Second)
			}

			for {
				select {
				case <-quit: // closed channel 不会阻塞，因此可用作退出通知。
					return
				default: // 执行正常任务。
					task()
				}
			}
		}(i)
	}

	time.Sleep(time.Second * 5) // 让测试 goroutine 运行一会
	close(quit)
	wg.Wait()
}
