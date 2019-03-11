// bubble sort 优化版本
package main

import "fmt"

func main() {
	array := []int{9, 8, 7, 6, 5, 4, 3, 2, 1, 11, 12, 13, 14, 15}
	fmt.Println("before sort array: ", array)
	arrayLen := len(array)
	flag := arrayLen
	for {
		if flag == 0 {
			break
		}
		flag = 0
		for i := 1; i < arrayLen; i++ {
			if array[i-1] > array[i] {
				tmp := array[i]
				array[i] = array[i-1]
				array[i-1] = tmp
				flag = i //没有任何交换操作,所以已经完全排序好了.这时flag = j这句不会执行,使得flag = 0.从而触发break
				// 减少了外层循环的次数。因为其实已经有序了，之后的外层循环可以不用再继续进行
			}
		}
	}
	fmt.Println("after sort array: ", array)
}
