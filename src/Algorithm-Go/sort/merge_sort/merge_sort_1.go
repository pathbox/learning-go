package main

import "fmt"

func mergeSort(r []int) []int {
	length := len(r)
	if length <= 1 { // 先递归按中间不断拆分，到数组r长度为1，然后递归出栈进行排序合并
		return r
	}
	mid := length / 2           // 当r是元素为2个的数组时，递归结束，开始将left和right进行merge排序
	left := mergeSort(r[:mid])  // 先不断的递归左部分，递归返回得到最大左部分
	right := mergeSort(r[mid:]) // 再不断的递归右部分，递归返回得到最大右部分
	return merge(left, right)
}

// 将两个有序数组，合并成一个有序数组
func merge(left, right []int) (result []int) {
	l, r := 0, 0
	for l < len(left) && r < len(right) { // 停止的条件其实是，left或right数组有一个遍历完
		if left[l] < right[r] { // left数组的值更小，取left的值
			result = append(result, left[l])
			l++
		} else { // right数组的值更小，取right的值
			result = append(result, right[r])
			r++
		}
	}
	// 剩下的数肯定都大于(或小于)result中的数
	result = append(result, left[l:]...)  // 如果 left还有没遍历的，全部加到result尾部
	result = append(result, right[r:]...) // 如果 right还有没遍历的，全部加到result尾部， 这两步只会有一步有append到result
	return
}

func main() {
	a := []int{9, 4, 2, 5, 6, 7, 1, 0, 3, 11, 75, 23}
	fmt.Println("Before Sort: ", a)
	result := mergeSort(a)
	fmt.Println("After Sort:  ", result)
}
