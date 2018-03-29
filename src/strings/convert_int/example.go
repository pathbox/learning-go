package main

import (
	"fmt"
	"strconv"
)

func main() {
	s := "123456"
	i, _ := strconv.Atoi(s)
	fmt.Println(i)
	//string到int64
	i64, _ := strconv.ParseInt(s, 10, 64)
	fmt.Println(i64)
	//int到string
	s = strconv.Itoa(i)
	fmt.Println(s)
	//int64到string
	ib := 100
	r := strconv.FormatInt(int64(ib), 10)
	fmt.Println(r)
	//string到float32(float64)
	float, _ := strconv.ParseFloat(s, 32/64)
	//float到string
	sa := strconv.FormatFloat(float, 'E', -1, 32)
	sb := strconv.FormatFloat(float, 'E', -1, 64)
	fmt.Println(sa, sb)
	// 'b' (-ddddp±ddd，二进制指数)
	// 'e' (-d.dddde±dd，十进制指数)
	// 'E' (-d.ddddE±dd，十进制指数)
	// 'f' (-ddd.dddd，没有指数)
	// 'g' ('e':大指数，'f':其它情况)
	// 'G' ('E':大指数，'f':其它情况)
}
