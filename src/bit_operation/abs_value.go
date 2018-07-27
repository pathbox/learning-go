package main

import "fmt"

func main() {
	a := -1
	fmt.Println(a)
	r := abs(a)
	fmt.Println(r)
}

func abs(a int) int {
	r := a >> 32
	return a ^ r - r
}
