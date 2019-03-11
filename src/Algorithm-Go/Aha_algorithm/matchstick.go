package main

import (
	"fmt"
)

// 输入的值 需要多少根火柴棍
func numCount(x int) int {
	num := 0                                 // 用来计数的变量
	f := []int{6, 2, 5, 5, 4, 5, 6, 3, 7, 6} // 用一个数组来记录0-9个数字需要用多少根火柴棍
	for x/10 != 0 {                          // 计算每一数字位，需要的火柴棍数量
		num = num + f[x%10]
		x = x / 10 // 去掉 x的末尾数字，例如x的值为123则现在x的值为12
	}

	// 最后加上此时x（最高位）所需要的火柴棍的根数
	num = num + f[x]
	return num // 返回需要火柴棍的总根数
}

func main() {
	a, b, c, m, sum := 0, 0, 0, 0, 0
	m = 24
	limit := 1111
	for a = 0; a <= limit; a++ {
		for b = 0; b <= limit; b++ {
			c = a + b
			if numCount(a)+numCount(b)+numCount(c) == m-4 { // + 和等号占用 4个火柴棍
				fmt.Printf("%d+%d=%d\n", a, b, c)
				sum = sum + 1
			}
		}
	}
	fmt.Printf("一共可以评出%d个不同的等式\n", sum)
}

/*
一共有24个火柴棍，进行a+b=c加法拼接。 可以拼接出多少种加法运算
*/
