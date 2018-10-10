package main

import "fmt"

func main() {
	a := 200
	fmt.Printf("%d\n", a>>1)
	fmt.Printf("%d\n", a<<1)
}

// 整数运算，左移1位相当于乘以2，右移1位相当于除以2
