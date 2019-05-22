package main

import (
	"fmt"
	"reflect"
)

func main() {
	s1 := S{}
	s2 := &S{}
	s1.Name()
	s1.Age()

	s2.Name()
	s2.Age()

	t1 := reflect.TypeOf(s1)
	fmt.Println("===============t1")
	for i := 0; i < t1.NumMethod(); i++ {
		fmt.Println(t1.Method(i).Name)
	}

	t2 := reflect.TypeOf(s2)
	fmt.Println("===============t2")
	for i := 0; i < t2.NumMethod(); i++ {
		fmt.Println(t2.Method(i).Name)
	}

	t3 := t2.Elem()
	fmt.Println("===============t3")
	for i := 0; i < t3.NumMethod(); i++ {
		fmt.Println(t3.Method(i).Name)
	}
}

type S struct {
}

func (s S) Name() {
	fmt.Println("Name")
}

func (s *S) Age() {
	fmt.Println("Age")
}
