package main

import (
	"fmt"
)

func main() {
	fmt.Println("c return:", *(c())) // 打印结果为 c return: 2
}

func c() *int {
	var i int
	defer func() {
		i++
		fmt.Println("c defer2:", i) // 打印结果为 c defer: 2
	}()
	defer func() {
		i++
		fmt.Println("c defer1:", i) // 打印结果为 c defer: 1
	}()
	return &i
}
