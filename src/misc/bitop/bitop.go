package main

import (
	"fmt"
)

func main() {
	fmt.Printf("%b\n", 0110&1011)  // AND, 0010
	fmt.Printf("%b\n", 0110|1011)  // OR, 1111
	fmt.Printf("%b\n", 0110^1011)  // XOR, 1101
	fmt.Printf("%b\n", 0110&^1011) // AND NOT, 清除标志位

	a := 0
	fmt.Printf("%b\n", a)
	a |= 1 << 2 //0000100: 在 bit2 设置标志位
	fmt.Printf("%b\n", a)
	a |= 1 << 6 // 1000100: 在 bit6 设置标志位
	fmt.Printf("%b\n", a)
	a = a &^ (1 << 6) // 0000100: 清除 bit6 标志位
	fmt.Printf("%b\n", a)

	// 取反用"^"
	x := 1
	fmt.Printf("%b, %b\n", x, ^x)
}
