package main

import "fmt"

func test(x [2]int) {
	fmt.Printf("x: %p\n", &x)
	x[1] = 1000
}

func main() {
	a := [2]int{}
	fmt.Printf("a: %p\n", &a)
	test(a)
	fmt.Println(a)
}
