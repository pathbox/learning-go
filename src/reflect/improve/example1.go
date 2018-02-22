package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type TestObj struct {
	field1 string
}

func main() {
	struct_ := &TestObj{}                                            //初始化一个TestObj{}
	field, _ := reflect.TypeOf(struct_).Elem().FieldByName("field1") // 获取field1 字段
	fmt.Println(field)
	field1Ptr := uintptr(unsafe.Pointer(struct_)) + field.Offset // 获取field字段的指针
	*((*string)(unsafe.Pointer(field1Ptr))) = "hello"            // 给field1Ptr指针所指的地址赋值 hello
	fmt.Println(struct_)                                         // 赋值后的 TestObj{}
	fmt.Println(struct_.field1)                                  // hello  这时候 field1的值就为hello
}

// fieldPtr := uintptr(structPtr) + field.Offset
