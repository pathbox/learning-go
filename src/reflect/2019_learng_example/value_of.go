package main

import (
	"fmt"
	"reflect"
)

type ss struct {
	A int
	a int
}

func (s ss) Method1(i int) string  { return "结构体方法1" }
func (s *ss) Method2(i int) string { return "结构体方法2" }

func main() {
	v1 := reflect.ValueOf(ss{})                   // 结构体
	v2 := reflect.ValueOf(&ss{})                  // 结构体指针
	v3 := reflect.ValueOf(&ss{}).Elem()           // 可寻址结构体
	v4 := reflect.ValueOf(&ss{}).Elem().Field(0)  // 可寻址结构体的共有字段
	v5 := reflect.ValueOf(&ss{}).Elem().Field(1)  // 可寻址结构体的私有字段
	v6 := reflect.ValueOf(&ss{}).Method(0)        // 结构体指针的方法
	v7 := reflect.ValueOf(&ss{}).Elem().Method(0) // 结构体的方法

	fmt.Println(v1.CanAddr()) // false
	fmt.Println(v2.CanAddr()) // false
	fmt.Println(v3.CanAddr()) // true
	fmt.Println(v4.CanAddr()) // true
	fmt.Println(v5.CanAddr()) // true
	fmt.Println(v6.CanAddr()) // false
	fmt.Println(v7.CanAddr()) // false
	fmt.Println("----------")
	fmt.Println(v1.CanSet()) // false
	fmt.Println(v2.CanSet()) // false
	fmt.Println(v3.CanSet()) // true
	fmt.Println(v4.CanSet()) // true
	fmt.Println(v5.CanSet()) // false
	fmt.Println(v6.CanSet()) // false
	fmt.Println(v6.CanSet()) // false
	fmt.Println("----------")
	fmt.Println(v1.CanInterface()) // true
	fmt.Println(v2.CanInterface()) // true
	fmt.Println(v3.CanInterface()) // true
	fmt.Println(v4.CanInterface()) // true
	fmt.Println(v5.CanInterface()) // false
	fmt.Println(v6.CanInterface()) // true
	fmt.Println(v7.CanInterface()) // true
}
