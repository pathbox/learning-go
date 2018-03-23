// 无名返回值的情况

package main

import (
	"fmt"
)

func main() {
	fmt.Println("return:", a()) // 打印结果为 return: 0
}

func a() int {
	var i int
	defer func() {
		i++
		fmt.Println("defer2:", i) // 打印结果为 defer: 2
	}()
	defer func() {
		i++
		fmt.Println("defer1:", i) // 打印结果为 defer: 1
	}()
	return i
}
