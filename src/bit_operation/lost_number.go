package main

import "fmt"

func main() {
	ary := []int{1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 7, 7, 8, 8, 9, 9}
	var lostNum int
	for _, item := range ary {
		lostNum ^= item
	}
	fmt.Printf("缺失的数字为:%d\n", lostNum)
}

// 自己与自己异或结果为0，2.异或满足交换律。因此我们将这些数字全异或一遍，结果就一定是那个仅出现一个的那个数
