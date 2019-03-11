// 简单递归法

// F(1) = 1;
// F(2) = 2;
// F(n) = F(n-1)+F(n-2)（n>=3）

// 时间复杂度： O(2^N)

package main

import (
	"fmt"
)

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

	return getClimbingWays(n-1) + getClimbingWays(n-2)
}
