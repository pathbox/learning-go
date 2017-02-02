package main

import (
	"fmt"
)

func adder() func(int) int {
	sum := 0
	innerfunc := func(x int) int {
		sum += x
		return sum
	}
	return innerfunc
}

func main() {
	pos, neg := adder(), adder()
	for i := 0; i < 100; i++ {
		fmt.Println(pos(i), neg(-2*i))
	}
}
