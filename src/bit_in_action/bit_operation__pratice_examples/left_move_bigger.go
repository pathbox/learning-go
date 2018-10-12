package main

import "fmt"

func main() {
	var a int8 = 3
	fmt.Printf("%08b\n", a)
	fmt.Printf("%08b\n", a<<1)
	fmt.Printf("%08b\n", a<<2)
	fmt.Printf("%08b\n", a<<3)
}

// 左移会按位移增大
