package main

import (
	"fmt"
)

// 被关闭的信道会禁止数据流入, 是只读的。我们仍然可以从关闭的信道中取出数据，但是不能再写入数据了。

func main() {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3

	// 显式地关闭信道
	close(ch)

	for v := range ch {
		fmt.Println(v)
	}
}
