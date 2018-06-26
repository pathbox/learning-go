package main

import (
	"fmt"
	"reflect"
)

func main() {
	a := []int{1, 2, 3}
	b := []int{1, 2, 3}

	fmt.Printf("%p,%p\n", a, b)

	r := reflect.DeepEqual(a, b)
	fmt.Println("a == b?:", r)

	m := make(map[string]string)
	n := make(map[string]string)

	m["name"] = "nice"
	n["name"] = "nice"
	// n["name"] = "nice1"

	fmt.Printf("%p,%p\n", m, n)

	rr := reflect.DeepEqual(m, n)

	fmt.Println("m == n?:", rr)
}
