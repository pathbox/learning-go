package main

import (
	"fmt"
)

func main() {
	array := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} // 有序的数组
	target := 6

	binarySearch(target, array)
}

func binarySearch(target int, array []int) {
	fmt.Println("My target is: ", target)

	lo, hi := 0, len(array)-1
	for lo <= hi { // 非递归二分法
		mid := (lo + hi) >> 1
		if array[mid] < target {
			lo = mid + 1
		} else if array[mid] > target {
			hi = mid - 1
		} else {
			fmt.Printf("target:%d is here", target)
			return
		}
	}
	fmt.Printf("target:%d is not here", target)
}
