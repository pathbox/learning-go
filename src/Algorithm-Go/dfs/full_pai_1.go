package main

import "fmt"

var ary = []int{1, 2, 3}
var theLen = len(ary)

func main() {

	for i, val := range ary {
		var out []int
		book := make([]int, theLen, theLen) // 不适合设为全局变量
		dfsFull(val, i, out, book)
	}

}

func dfsFull(item, index int, out, book []int) {
	out = append(out, item)
	book[index] = 1

	if len(out) == theLen {
		for _, v := range out {
			fmt.Printf("%d", v)
		}
		out = []int{}
		fmt.Println()
		return
	}

	for i, val := range ary { // 循环item的邻结点
		if book[i] == 0 {
			dfsFull(val, i, out, book)
			book[i] = 0
		}
	}
}
