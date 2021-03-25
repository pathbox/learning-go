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
*/
