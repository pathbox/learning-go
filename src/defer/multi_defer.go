package main

import (
	"fmt"
)

// 多个 defer 注册，按 FILO 次序执行。哪怕函数或某个延迟调用发生错误，这些调用依旧会被执行。  先进后出

func test(x int) {
	defer fmt.Println("a")
	defer fmt.Println("b")

	defer func() {
		fmt.Println(100 / x) // div0 异常未被捕获，逐步往外传递，最终终⽌止进程。
	}()

	defer fmt.Println("c")
}

func main() {
	test(0)
}
