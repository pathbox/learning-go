package main

import (
	"fmt"
)

var a = [10]int{}

var n int

var book = [10]int{}

func dfs(step int) { // step 表示现在站在第几个盒子前面， 或者 现在要拍第几位的数字

	if step == n+1 { // 初始索引为0，如果站在第n+1个盒子面前，则表示前n个盒子已经放好扑克牌了
		for _, v := range a {
			if v != 0 {
				fmt.Print(v)
			}
		}
		fmt.Println()
		return // 返回之前的一步（最近一次调用dfs函数的地方）
	}

	// 此时站在第step个盒子面前，应该放哪一张牌呢？
	// 按照1、2、3 ... n的顺序一一尝试
	for i := 1; i <= n; i++ {
		// 判断扑克牌i是否还在手上
		if book[i] == 0 { // book[i]等于0表示i号扑克牌在手上
			// 开始尝试使用扑克牌i
			a[step] = i // 将i号扑克牌放入到第step个盒子中
			book[i] = 1 // 将book[i]设为1，表示i号扑克牌已经不在手上

			// 第step个盒子（全排列的第step位）已经放好扑克牌，接下来需要走到下一个盒子面前
			dfs(step + 1) // 这里通过函数的递归调用来实现
			book[i] = 0   // 将刚才尝试的扑克牌收回，才能进行下一次尝试

		}
	}
	return
}

func main() {
	n = 3
	for i := 1; i <= n; i++ {
		a[i] = 0
		book[i] = 0
	}
	dfs(1) // 首先站在1号盒子面前
}
