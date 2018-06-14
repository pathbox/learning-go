package main

import (
	"bytes"
	"fmt"
)

func main() {
	foo := []byte("foo")
	bar := []byte("bar")
	fmt.Println("foo bar: ", foo, bar)
	fmt.Println(bytes.Contains([]byte("seafood"), foo))
	fmt.Println(bytes.Contains([]byte("seafood"), bar))
	fmt.Println(bytes.Contains([]byte("seafood"), []byte("")))
	fmt.Println(bytes.Contains([]byte("北京"), []byte("北京")))
}
