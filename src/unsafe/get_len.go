package main

import (
	"fmt"
	"unsafe"
)

func main() {
	s := make([]int, 9, 20)
	l :=*(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&s)) + uintptr(8)))

	fmt.Println(l, len(s)) // 9 9

	c := *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&s))+uintptr(16)))

	fmt.Println(c, cap(s)) // 20 20
}

/*
Len: &s => pointer => uintptr => pointer => *int => int
Cap: &s => pointer => uintptr => pointer => *int => int
Pointer先转为 uintptr, uintptr进行计算之后, 需要再转为Pointer, 然后才能被 (*int)转换，之后加 * 号取到指针的值
*/
