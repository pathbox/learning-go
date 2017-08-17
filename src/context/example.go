package main

import (
	"context"
	"fmt"
	"time"
)

func childFunc(cont context.Context, num *int) {
	ctx, _ := context.WithCancel(cont) // 创建子 context
	for {
		select {
		case <-ctx.Done():
			fmt.Println("child Done: ", ctx.Err())
			return
		}
	}
}

func main() {
	gen := func(ctx context.Context) <-chan int {
		dst := make(chan int)
		n := 1
		go func() {
			for {
				select {
				case <-ctx.Done():
					fmt.Println("parent Done : ", ctx.Err())
					return // returning not to leak the goroutine
				case dst <- n:
					n++
					go childFunc(ctx, &n)
				}
			}
		}()
		return dst //  dst channel 会 进行 <- 取操作，取出channel中的值，然后返回
	}
	ctx, cancel := context.WithCancel(context.Background())
	for n := range gen(ctx) {
		fmt.Println(n)
		if n >= 5 {
			break
		}
	}
	cancel()
	time.Sleep(2 * time.Second)
}

/*

context 对象形成一棵树：当一个 Context 对象被取消时，继承自它的所有 Context 都会被取消。

在上面的例子中，主要描述的是通过一个channel实现一个为循环次数为5的循环，
在每一个循环中产生一个goroutine，每一个goroutine中都传入context，在每个goroutine中通过传入ctx创建一个子Context,并且通过select一直监控该Context的运行情况，当在父Context退出的时候，代码中并没有明显调用子Context的Cancel函数，但是分析结果，子Context还是被正确合理的关闭了，这是因为，所有基于这个Context或者衍生的子Context都会收到通知，这时就可以进行清理操作了，最终释放goroutine，这就优雅的解决了goroutine启动后不可控的问题。

3.5 Context 使用原则
不要把Context放在结构体中，要以参数的方式传递
以Context作为参数的函数方法，应该把Context作为第一个参数，放在第一位。
给一个函数方法传递Context的时候，不要传递nil，如果不知道传递什么，就使用context.TODO
Context的Value相关方法应该传递必须的数据，不要什么数据都使用这个传递
Context是线程安全的，可以放心的在多个goroutine中传递
*/
