package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	const total = 1000
	var wg sync.WaitGroup
	wg.Add(total)
	now := time.Now()
	for i := 0; i < total; i++ {
		go func() {
			defer wg.Done()
			requestWork(context.Background(), "any")
		}()
	}
	wg.Wait()
	fmt.Println("elapsed:", time.Since(now))
	time.Sleep(time.Minute * 2)
	fmt.Println("number of goroutines:", runtime.NumGoroutine()) // 打印正在运行的goroutine数量
}

func requestWork(ctx context.Context, job interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	done := make(chan error) // fix: done := make(chan error, 1)
	go func() {
		done <- hardWork(job)
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func hardWork(job interface{}) error {
	time.Sleep(time.Minute)
	return nil
}

// https://gocn.vip/topics/11866
/*
requestWork 所在的goroutine被ctx超时控制了，但是requestWork内部的goroutine没有
goroutine 泄露了，让我们看看为啥会这样呢？首先，requestWork 函数在 2 秒钟超时后就退出了，一旦 requestWork 函数退出，那么 done channel 就没有 goroutine 接收了，等到执行 done <- hardWork(job) 这行代码的时候就会一直卡着写不进去，导致每个超时的请求都会一直占用掉一个 goroutine，这是一个很大的 bug，等到资源耗尽的时候整个服务就失去响应了。

那么怎么 fix 呢？其实也很简单，只要 make chan 的时候把 buffer size 设为 1，如下：

done := make(chan error, 1) goroutine不会再阻塞，channel多了些内存占用
这样就可以让 done <- hardWork(job) 不管在是否超时都能写入而不卡住 goroutine。此时可能有人会问如果这时写入一个已经没 goroutine 接收的 channel 会不会有问题，在 Go 里面 channel 不像我们常见的文件描述符一样，不是必须关闭的，只是个对象而已，close(channel) 只是用来告诉接收者没有东西要写了，没有其它用途
*/
