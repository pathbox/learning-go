package main

import (
	"fmt"
)

func test() {
	x, y := 10, 20

	defer func(i int) { // i is (x)
		fmt.Println("defer: ", i, y)
	}(x)

	x += 10
	y += 100
	fmt.Println("x = ", x, "y = ", y)
}

func main() {
	test()
}
