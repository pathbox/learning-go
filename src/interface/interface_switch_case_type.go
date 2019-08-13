package main

import (
	"fmt"
)

type User struct {
	Name string
}

func main() {
	any := User{
		Name: "fidding",
	}
	myType(any)
	any2 := "fidding"
	myType(any2)
	any3 := int32(123)
	myType(any3)
	any4 := int64(123)
	myType(any4)
	any5 := []int{1, 2, 3, 4, 5}
	myType(any5)
}

func myType(value interface{}) {
	switch value.(type) {
	case string:
		// 将interface转为string字符串类型
		op, ok := value.(string)
		fmt.Println(op, ok)
	case int32:
		// 将interface转为int32类型
		op, ok := value.(int32)
		fmt.Println(op, ok)
	case int64:
		// 将interface转为int64类型
		op, ok := value.(int64)
		fmt.Println(op, ok)
	case User:
		// 将interface转为User struct类型，并使用其Name对象
		op, ok := value.(User)
		fmt.Println(op.Name, ok)
	case []int:
		// 将interface转为切片类型
		op := make([]int, 0)
		op = value.([]int)
		fmt.Println(op)
	case []string:
		op := make([]string, 0)
		op = value.([]string)
		fmt.Println(op)
	case []interface{}:
		op := make([]interface{}, 0)
		op = value.([]interface{})
		fmt.Println(op)
	default:
		fmt.Println("unknown")

	}
}

// 解析具体的interface{} 为什么值，转换该值