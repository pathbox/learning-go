package main

import (
	"fmt"
)

func main() {
	var a int
	a = 1
	a.aaaa()
	a.bbbb()
}

func (i *int) aaaa() {
	fmt.Println(i)
}

func (i int) bbbb() {
	fmt.Println(i)
}

// ./example4.go:10: a.aaaa undefined (type int has no field or method aaaa)
// ./example4.go:11: a.bbbb undefined (type int has no field or method bbbb)
// ./example4.go:14: cannot define new methods on non-local type int
// ./example4.go:18: cannot define new methods on non-local type int
