package main

import (
	"fmt"
)

func fb(n int) int {
	if n == 1 {
		return 1
	}

	if n == 2 {
		return 1
	}

	if n < 0 {
		return 0
	}

	s := fb(n-1) + fb(n-2)
	return s
}

func main() {
	n := 7
	r := fb(n)
	fmt.Println(r)

}
