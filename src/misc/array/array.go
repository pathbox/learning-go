package main

import "fmt"

func main() {
	var numbers [5]int
	numbers[0] = 2
	numbers[3] = numbers[0] - 3
	numbers[1] = numbers[2] + 5
	fmt.Println(numbers)
}
