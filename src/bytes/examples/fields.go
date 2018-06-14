package main

import (
	"bytes"
	"fmt"
)

func main() {
	a := bytes.Fields([]byte("  foo bar  baz   ")) // return [][]byte
	for i, b := range a {
		fmt.Println(i, "-", string(b))
	}
	fmt.Printf("Fields are: %q", bytes.Fields([]byte("  foo bar  baz   ")))
}
