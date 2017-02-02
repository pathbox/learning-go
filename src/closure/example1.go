package main

import (
	"fmt"
)

func adder() func(int) int {
	sum := 0 // 只会执行一次，之后返回的是一个闭包，闭包里面的代码执行循环了100次
	innerfunc := func(x int) int {
		sum += x
		return sum
	}
	return innerfunc
}

func main() {
	pos, neg := adder(), adder()
	for i := 0; i < 100; i++ {
		fmt.Println(pos(i), neg(-2*i))
	}
}
