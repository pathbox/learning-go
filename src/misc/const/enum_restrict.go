package main

import (
	"fmt"
)

type Color int

const (
	Black Color = iota
	Red
	Blue
)

func test(c Color) {
	fmt.Println(c)
}

func main() {
	c := Black
	test(c)

	x := 1
	//test(x) // Error: cannot use x (type int) as type Color in function argument
	_ = x

	test(1) // 常量会被编译器自动转换。

	d := Blue
	test(d)
}
