package main

import "fmt"

func main() {

	ary := [6]int{1, 2, 3, 4, 5, 6}

	fmt.Println(ary)
	s1 := ary[0:3]
	fmt.Println(s1)

	ary[1] = 100

	fmt.Println(s1)

	s1 = []int{11, 12, 13, 14, 15, 16, 18, 19}
	fmt.Println(s1)

	fmt.Println(ary)

}
