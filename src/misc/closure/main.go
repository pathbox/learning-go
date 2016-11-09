package main

import "fmt"

func nextNum() func() int {
	i, j := 2, 2
	return func() int {
		var tmp = i + j
		i, j = j, tmp
		return tmp
	}
}

func main() {
	nextNumFunc := nextNum()
	for i := 0; i < 10; i++ {
		fmt.Println(nextNumFunc())
	}
}
