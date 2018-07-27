package main

import "fmt"

func main() {
	a := int64(-11)
	fmt.Println(a)
	r := abs(a)
	fmt.Println(r)
}

func abs(a int64) int64 {
	r := a >> 31

	fmt.Println("r:", r)
	// return a ^ r - r
	if r == 0 { // 说明 a是正整数，直接返回
		return a
	} else {
		return ^a + 1 // 说明a是负整数，取相反符号,即求正整数：反码+1
	}
}
