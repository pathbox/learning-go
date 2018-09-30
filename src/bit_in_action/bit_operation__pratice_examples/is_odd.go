package main

import "fmt"

// 判断奇偶数， & 1 得到 0 就是偶数， 1 就是奇数
func main() {
	for i := 0; i <= 100; i++ {
		if i&1 == 1 {
			fmt.Printf("%d is 奇数\n", i)
		} else {
			fmt.Printf("%d is 偶数\n", i)
		}
	}
}
