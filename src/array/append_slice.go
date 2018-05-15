package main

import (
	"fmt"
)

func main() {
	s := []int{1, 2, 3}
	ss := s[1:]
	fmt.Println(ss)
	for i := range ss {
		ss[i] += 10
	}
	fmt.Println(s)
	ss = append(ss, 4) // ss扩容后，底层指向新的数据而不是s
	fmt.Println(ss)
	for i := range ss {
		ss[i] += 10
	}
	fmt.Println(s)
	fmt.Println(ss)
}
