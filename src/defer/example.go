package main

import "fmt"

func main() {

	a := 1
	b := 2

	defer put(a)
	defer put(b)
}

func put(i int) {
	fmt.Println(i)
}

// defer is a stack  先进后出
