package main

import (
	"fmt"
)

type A struct {
	User string
}

type B struct {
	A
}

func main() {
	var p B
	p.User = "Hello"
	// wrong
	// p := B {
	// User: "hello"
	//}

	fmt.Println(p)
}
