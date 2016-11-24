package main

import (
	"fmt"
)

type A struct {
	a int
}

func (s A) f() {
	fmt.Println("A f()")
}

type B struct {
	A
	b int
}

func main() {
	s := A{}
	s.f()
	s1 := B{}
	s1.f()
	s1.A.f()
}
