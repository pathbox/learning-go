package main

import (
	"fmt"
	"time"
)

func main() {
	go foo()

	fmt.Println("Here is main goroutinue") // 主 goroutinue 执行完之后没有等待子 goroutinue，子 goroutinue 会结束被垃圾回收
}

func foo() {
	time.Sleep(time.Second * 3)
	fmt.Println("Here is sub goroutinue")
}
