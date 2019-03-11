package main

import (
	"fmt"
)

func main() {

	fmt.Println("Start bulle sort")
	array := [10]int{1, 3, 5, 7, 9, 2, 4, 6, 8, 10}

	array_len := len(array)

	for i := 1; i < array_len; i++ { // 9次 len-1  i 从1开始，只是为了 👇的array_len - i
		for j := 0; j < array_len-i; j++ { // j j+1的值范围是[0-9] j从0开始进行，循环 len-1次
			if array[j] > array[j+1] {
				array[j], array[j+1] = array[j+1], array[j]
			}
			fmt.Println(j, "-", array)
		}
	}

	fmt.Println("After bubble sort")
	fmt.Println(array)
}
