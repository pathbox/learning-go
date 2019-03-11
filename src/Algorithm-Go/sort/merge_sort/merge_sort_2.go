package main

import "fmt"

func main() {
	a := []int{9, 4, 2, 5, 6, 7, 1, 0, 3, 11, 75, 23}
	fmt.Println("Before Sort: ", a)
	result := mergeSort(a)
	fmt.Println("After Sort:  ", result)
}

func mergeSort(array []int) []int {
	if len(array) < 2 {
		return array
	}

	mid := len(array) / 2
	return merge(mergeSort(array[:mid]), mergeSort(array[mid:]))
}

func merge(left, right []int) []int {
	size, l, r := len(left)+len(right), 0, 0
	slice := make([]int, size, size)

	for i := 0; i < size; i++ {
		if l > len(left)-1 && r <= len(right)-1 {
			slice[i] = right[r]
			r++
		} else if r > len(right)-1 && l <= len(left)-1 {
			slice[i] = left[l]
			l++
		} else if left[l] < right[r] {
			slice[i] = left[l]
			l++
		} else {
			slice[i] = right[r]
			r++
		}
	}
	return slice
}
