/*
# 算法分析
# 1．先从数列中取出一个数作为基准数。
# 2．分区过程，将比这个数大的数全放到它的右边，小于或等于它的数全放到它的左边。
# 3．再对左右区间重复第二步，直到各区间只有一个数。(终止条件是:直到各个区间只有一个数)
在最优的情况下，快速排序算法的时间复杂度为O(nlogn). 最坏情况，需要进行n‐1递归调用,时间复杂度为O(n^2)
*/

package main

import (
	"fmt"
)

var array []int

func main() {
	array = []int{9, 8, 7, 6, 5, 4, 3, 2, 1, 11, 12, 13, 14, 15}
	fmt.Println("before sort array: ", array)
	arrayLen := len(array)
	l, r := 0, arrayLen-1
	quickSort(l, r)

	fmt.Println("after sort array: ", array)

}

func quickSort(left, right int) {
	var i, j, temp int
	if left > right {
		return
	}

	i = left
	j = right
	temp = array[left]

	for i != j {
		for array[j] >= temp && i < j { //顺序很重要，要先从右往左找, 找到比基准数小的数停下
			j--
		}
		for array[i] <= temp && i < j { //再从左往右找，找到比基准数大的数停下
			i++
		}
		if i < j { //当哨兵i和哨兵j没有相遇时 交换i和j
			array[i], array[j] = array[j], array[i]
		}
	}
	//i==j 相遇了 最终将基准数归位
	array[left], array[i] = array[i], array[left]
	// 递归 分而治之
	quickSort(left, i-1) //继续处理左边的，这里是一个递归的过程

	quickSort(i+1, right) //继续处理右边的，这里是一个递归的过程

}
