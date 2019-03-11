package main

import (
	"fmt"
)

func main() {
	array := [100]int{6, 3, 1, 7, 5, 8, 9, 2, 4}
	head, tail := 0, 9

	for head < tail {
		fmt.Print(array[head])
		head++
		array[tail] = array[head]
		tail++
		head++
	}
}
