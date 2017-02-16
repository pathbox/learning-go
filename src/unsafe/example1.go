package main

import (
	"fmt"
	"unsafe"
)

func main() {
	var n int64 = 5
	var pn = &n
	var pf = (*float64)(unsafe.Pointer(pn))

	// now pn and pf are pointing at same memory address
	fmt.Println(*pf)
	fmt.Println(*pn)
	fmt.Println(pf)
	fmt.Println(pn)
	*pf = 3.14159
	fmt.Println(n)
}
