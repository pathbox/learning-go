package main

import (
	"fmt"
)

type user struct {
	name string
	age  uint64
}

func bad() {
	u := []user{
		{"asong", 23},
		{"song", 19},
		{"asong2020", 18},
	}
	for _, v := range u {
		if v.age != 18 { // v是拷贝值，不会修改原来的结构体值
			v.age = 20
		}
	}
	fmt.Println(u)
}

// [{asong 23} {song 19} {asong2020 18}]

func good() {
	u := []user{
		{"asong", 23},
		{"song", 19},
		{"asong2020", 18},
	}
	for k, v := range u {
		if v.age != 18 {
			u[k].age = 18
		}
	}
	fmt.Println(u)
}
