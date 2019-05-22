package main

import (
	"fmt"
	"reflect"
)

type Foo struct {
	A int `tag1:"Tag1" tag2:"Second Tag"`
	B string
}

func main() {
	s := "String 字符串"
	fo := Foo{A: 10, B: "字段String字符串"}

	sVal := reflect.ValueOf(s)
	// 在没有获取指针的前提下，我们只能读取变量的值。
	fmt.Println(sVal.Interface())

	sPtr := reflect.ValueOf(&s)
	sPtr.Elem().Set(reflect.ValueOf("修改值1"))
	sPtr.Elem().SetString("修改值2")
	// 修改指针指向的值，原变量改变
	fmt.Println(s)
	fmt.Println(sPtr) // 要注意这是一个指针变量，其值是一个指针地址
	foType := reflect.TypeOf(fo)
	foKind := foType.Kind()
	fmt.Println("===", foType, foKind)
	foVal := reflect.New(foType)
	foVal.Elem().Field(0).SetInt(100) // 修改具体struct field的值，使用field index
	foVal.Elem().Field(1).SetString("B Nice")
	f2 := foVal.Elem().Interface().(Foo)

	fmt.Printf("%+v, %d, %s\n", f2, f2.A, f2.B)
}
