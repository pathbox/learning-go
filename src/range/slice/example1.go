package main

import (
	"fmt"
)

type user struct {
	name string
	age  uint64
}

// bad
func bad() {
	u := []user{
		{"asong", 23},
		{"song", 19},
		{"asong2020", 18},
	}
	n := make([]*user, 0, len(u))
	for _, v := range u {
		n = append(n, &v)
	}
	fmt.Println(n)
	for _, v := range n {
		fmt.Println(v)
	}
}

/*
[0xc00011a040 0xc00011a040 0xc00011a040]
&{asong2020 18}
&{asong2020 18}
&{asong2020 18}
*/

// good
func good() {
	u := []user{
		{"asong", 23},
		{"song", 19},
		{"asong2020", 18},
	}
	n := make([]*user, 0, len(u))
	for _, v := range u {
		o := v // here
		n = append(n, &o)
	}
	fmt.Println(n)
	for _, v := range n {
		fmt.Println(v)
	}
}

/*
在for range 中，变量v是用来保存迭代切片所得的值，因为v只被声明了一次，每次迭代的值都是赋值给v，该变量的内存地址始终未变，这样讲他的地址追加到新的切片中，该切片保存的都是同一个地址。
变量v的地址也并不是指向原来切片u[2]的，因为变量v的数据是切片拷贝的数据，是直接copy了结构体数据
*/
