package main

import (
	"fmt"
)

type FilterFunc func(ele interface{}) interface{}

/*
  公共操作:对数据进行特殊操作
*/
func Data(arr interface{}, filterFunc FilterFunc) interface{} {

	slice := make([]int, 0)
	array, _ := arr.([]int)

	for _, value := range array {

		integer, ok := filterFunc(value).(int)
		if ok {
			slice = append(slice, integer)
		}

	}
	return slice
}

/*
  具体操作:奇数变偶数（这里可以不使用接口类型,直接使用int类型)
*/
func EvenFilter(ele interface{}) interface{} {

	integer, ok := ele.(int)
	if ok {
		if integer%2 == 1 {
			integer = integer + 1
		}
	}
	return integer
}

/*
  具体操作:偶数变奇数（这里可以不使用接口类型,直接使用int类型)
*/
func OddFilter(ele interface{}) interface{} {

	integer, ok := ele.(int)

	if ok {
		if integer%2 != 1 {
			integer = integer + 1
		}
	}

	return integer
}

func main() {
	sliceEven := make([]int, 0)
	sliceEven = append(sliceEven, 1, 2, 3, 4, 5)
	sliceEven = Data(sliceEven, EvenFilter).([]int)
	fmt.Println(sliceEven) //[2 2 4 4 6]

	sliceOdd := make([]int, 0)
	sliceOdd = append(sliceOdd, 1, 2, 3, 4, 5)
	sliceOdd = Data(sliceOdd, OddFilter).([]int)
	fmt.Println(sliceOdd) //[1 3 3 5 5]

}
