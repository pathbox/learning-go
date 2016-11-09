package main

import (
	"fmt"
)

func main() {
L1:
	for x := 0; x < 3; x++ {
	L2:
		for y := 0; y < 5; y++ {
			if y > 2 {
				continue L2
			}
			if x > 1 {
				break L1
			}
			fmt.Println(x, ":", y, " ")
		}
		fmt.Println()
	}
}
