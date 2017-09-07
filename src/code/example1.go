// 使用两个 goroutine 交替打印序列，一个 goroutinue 打印数字， 另外一个goroutine打印字母， 最终效果如下 12AB34CD56EF78GH910IJ 。

package main

import (
	"fmt"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// make 三个channel
	chan_n := make(chan bool)
	chan_c := make(chan bool, 1) // 需要给chan 定义一个缓存容量空间，否则会死锁
	chan_done := make(chan bool)

	// 在两个goroutine中执行
	go printN(chan_c, chan_n)
	go printC(chan_c, chan_n, chan_done)

	chan_c <- true
	<-chan_done
}

// 打印数字
func printN(chan_c, chan_n chan bool) {

	for i := 1; i <= 10; i += 2 {
		<-chan_c // 阻塞等待
		for j := 0; j < 2; j++ {
			fmt.Print(i + j)
		}
		chan_n <- true // 为了printC 解锁
	}

}

// 打印字母
func printC(chan_c, chan_n, chan_done chan bool) {

	char_seq := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K"}
	for i := 0; i < 10; i += 2 {
		<-chan_n // 阻塞等待
		for j := 0; j < 2; j++ {
			fmt.Print(char_seq[i+j])
		}
		chan_c <- true // 为了 printN解锁
	}
	chan_done <- true
}

// 两个goroutine互相使用chan，互相阻塞等待。这种情况很容易产生死锁
