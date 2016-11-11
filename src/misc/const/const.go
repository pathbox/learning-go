package main

import (
	"fmt"
	"unsafe"
)

const x, y int = 1, 2
const s = "Hello World!"

const (
	a, b      = 10, 100
	c    bool = false
)

// 在常量组中，如不提供类型和初始化值，那么视作与上一常量相同。
const (
	s1 = "abc"
	s2 // s2 = "abc"
)

const (
	x1 = "abc"
	x2 = len(x1)
	x3 = unsafe.Sizeof(x2)
)

// 如果常量类型足以存储初始化值，那么不会引发溢出错误。
const (
	x4 byte = 100 // int to byte
	// x5 int  = 1e20 // float64 to int, overflows
)

func main() {
	const x = "xxxx" // 未使用局部常量不会引发编译错误
	fmt.Println(x)
	fmt.Println(s2)
	fmt.Println(x1)
	fmt.Println(x2)
	fmt.Println(x3)
	fmt.Println(x4)
	// fmt.Println(x5)
}
