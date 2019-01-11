package main

import "fmt"

func main() {
	table := make([][]bool, 5) // 5表示的是行数
	for i := range table {
		fmt.Println(i)
		table[i] = make([]bool, 9) // 9 表示的是列数
	}

	fmt.Println(table)
}
