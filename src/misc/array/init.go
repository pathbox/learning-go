package main

import (
	"fmt"
)

func main() {
	a := [3]int{1, 2}
	b := [...]int{1, 2, 3, 4}
	b1 := []int{1, 2, 3, 4}
	c := [5]int{2: 100, 4: 200}

	d := [...]struct {
		name string
		age  uint8
	}{
		{"user1", 10},
		{"user2", 20},
	}
	a2 := [2][3]int{{1, 2, 3}, {4, 5, 6}}
	b2 := [...][2]int{{1, 1}, {2, 2}, {3, 3}} // 第 2 纬度不能用 "..."。

	fmt.Println(a, b, b1, c, d, a2, b2)
}
