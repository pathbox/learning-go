package main

// func main() {
//     ch := func() <-chan int {
//         ch := make(chan int)
//         go func() {
//             for i := 0; ; i++ {
//                 ch <- i
//             }
//         } ()
//         return ch
//     }()

//     for v := range ch {
//         fmt.Println(v)
//         if v == 5 {
//             break  //下面的程序中后台Goroutine向管道输入自然数序列，main函数中输出序列。但是当break跳出for循环的时候，后台Goroutine就处于无法被回收的状态了
//         }
//     }
// }

// func main() {
//     ctx, cancel := context.WithCancel(context.Background())

//     ch := func(ctx context.Context) <-chan int {
//         ch := make(chan int)
//         go func() {
//             for i := 0; ; i++ {
//                 select {
//                 case <- ctx.Done():
//                     return
//                 case ch <- i:
//                 }
//             }
//         } ()
//         return ch
//     }(ctx)

// 		for v := range ch {
//         fmt.Println(v)
//         if v == 5 {
//             cancel()
//             break
//         }
//     }
// }

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*3)
	defer cancel()
	go task(ctx)
	time.Sleep(time.Second * 6)
}

func task(ctx context.Context) {
	ch := make(chan struct{}, 0)
	go func() {
		// 模拟4秒耗时任务
		time.Sleep(time.Second * 4)
		ch <- struct{}{}
	}()
	select {
	case <-ch:
		fmt.Println("done")
	case <-ctx.Done():
		fmt.Println("timeout")
	}
}
