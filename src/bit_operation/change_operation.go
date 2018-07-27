package main

import "fmt"

func main() {
	a := 11
	b := -22
	fmt.Println(a, b)
	a = signReversal(a)
	b = signReversal(b)

	fmt.Println(a, b)
}

func signReversal(a int) int {

	r := ^a

	return r + 1
}

/*
变换符号
变换符号就是正数变成负数，负数变成正数。

如对于-11和11，可以通过下面的变换方法将-11变成11

      1111 0101(二进制) –取反-> 0000 1010(二进制) –加1-> 0000 1011(二进制)

同样可以这样的将11变成-11

			0000 1011(二进制) –取反-> 1111 0100(二进制) –加1-> 1111 0101(二进制)
*/
