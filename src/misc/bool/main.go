package main

import "fmt"

// btoi returns 1 if b is true and 0 if false.

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

func itob(i int) bool {
	return i == 0
}

func main() {
	fmt.Println(btoi(true), btoi(false))
	fmt.Println(itob(1), itob(0))
}
