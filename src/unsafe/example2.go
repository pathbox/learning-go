package main

import (
	"fmt"
	"unsafe"
)

type Person struct {
	name   string
	age    int
	gender bool
}

func main() {
	a := [4]int{0, 1, 2, 3}
	p1 := unsafe.Pointer(&a[1])
	p3 := unsafe.Pointer(uintptr(p1) + 2*unsafe.Sizeof(a[0]))
	*(*int)(p3) = 6
	fmt.Println("a =", a) // a = [0 1 2 6]

	who := Person{"John", 30, true}
	pp := unsafe.Pointer(&who)
	fmt.Println(pp)
	pname := (*string)(unsafe.Pointer(uintptr(pp) + unsafe.Offsetof(who.name)))
	page := (*int)(unsafe.Pointer(uintptr(pp) + unsafe.Offsetof(who.age)))
	pgender := (*bool)(unsafe.Pointer(uintptr(pp) + unsafe.Offsetof(who.gender)))
	fmt.Println("pname", pname)
	fmt.Println("page", page)
	fmt.Println("pgender", pgender)
	*pname = "Alice"
	*page = 28
	*pgender = false
	pt := unsafe.Pointer(&who)
	fmt.Println(pt)
	fmt.Println(who) // {Alice 28 false}
}
