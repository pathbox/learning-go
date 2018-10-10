package main

import "fmt"

func main() {
	var a uint8 = 8
	fmt.Printf("%08b\n", a)
	b := (1 << 5)
	fmt.Printf("%08b\n", b)
	a = a | (1 << 5) // 将a的从右往左数第5+1位置为1
	fmt.Printf("%08b\n", a)

	var aa int8 = 13
	fmt.Printf("%08b\n", aa)
	aa = aa &^ (1 << 2)
	fmt.Printf("%08b\n", aa)
}

// 位移运算符提供了有趣的方式处理二进制值中特定位置的值。例如，下列的代码中，| 和 << 用于设置变量 a 的第6个 bit 位
// 使用 &^ 和位移运算符，我们可以取消设置一个值的某个位。例如，下面的示例将变量 aa 的第三位置为 0
