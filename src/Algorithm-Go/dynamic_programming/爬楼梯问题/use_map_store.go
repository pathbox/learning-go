// 使用了map存储已经计算过的值
// 时间复杂度和空间复杂度都是 N
package main

import (
	"fmt"
)

var storeMap = map[int]int{}

func main() {
	n := 10
	result := getClimbingWays(n)
	fmt.Println("The result is: ", result)
}

func getClimbingWays(n int) int {
	if n < 1 {
		return 0
	} else if n == 1 {
		return 1
	} else if n == 2 {
		return 2
	}

	if value, ok := storeMap[n]; ok {
		return value
	} else {
		value := getClimbingWays(n-1) + getClimbingWays(n-2)
		storeMap[n] = value
		return value
	}
}
