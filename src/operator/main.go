package main

import (
	"fmt"
)

// 获取0-n之间的所有偶数
func even(a int) (array []int) {
	for i := 0; i < a; i++ {
		if i&1 == 0 { // 位操作符&与C语言中使用方式一样
			array = append(array, i)
		}
	}
	return array
}

//互换两个变量的值
// 不需要使用第三个变量做中间变量
func swap(a, b int) (int, int) {
	a ^= b // 异或等于运算
	b ^= a
	a ^= b
	return a, b
}

// 左移、有移运算
func shifting(a int) int {
	a = a << 1
	a = a >> 1
	return a
}

// 变换符号
func nagation(a int) int {
	return ^a + 1
}

func main() {
	fmt.Printf("even: %v\n", even(100))
	a, b := swap(100, 200)
	fmt.Printf("swap: %d\t%d\n", a, b)
	fmt.Printf("shifting: %d\n", shifting(100))
	fmt.Printf("nagation: %d\n", nagation(100))
}
