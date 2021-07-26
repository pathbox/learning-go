package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ch := make(chan int)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case i := <-ch:
				fmt.Println(i)
			case <-time.After(3 * time.Second): // goroutine内部3秒超时
				fmt.Println("goroutine 1 timeout")
				return
			}
		}
	}()

	ch <- 1
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	background := context.Background()
	ctx, _ := context.WithTimeout(background, 10*time.Second) // 定义ctx 十秒超时
	wg.Add(1)

	go func(ctx context.Context) {
		defer wg.Done()
		for {
			select {
			case <-ticker.C:
				fmt.Println("tick")
			case <-time.After(5 * time.Second):
				fmt.Println("goroutine2 timed out")
			//上下文传递超时信息，结束goroutine
			case <-ctx.Done(): // 10秒之后，ctx.Done() 传递值，此时进行下面的return，结束该goroutine。 这其实是借助ctx，进行超时结束
				fmt.Println("goroutine2 done")
				return
			}
		}
	}(ctx)

	wg.Wait()
}
