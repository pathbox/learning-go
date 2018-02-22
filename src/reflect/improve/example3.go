package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type sliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

func main() {
	slice := []string{"hello", "world"}
	header := (*sliceHeader)(unsafe.Pointer(&slice))
	fmt.Println(header.Len)
	elementType := reflect.TypeOf(slice).Elem()
	secondElementPtr := uintptr(header.Data) + elementType.Size()
	firstElementPtr := uintptr(header.Data)
	*((*string)(unsafe.Pointer(secondElementPtr))) = "try"
	*((*string)(unsafe.Pointer(firstElementPtr))) = "Nice"
	fmt.Println(slice)
}

// slice 的秘诀在于取出指向数组头部的指针，然后具体的元素，通过偏移量来计算
