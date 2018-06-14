package main

import (
	"bytes"
	"fmt"
)

func main() {
	rs := bytes.Runes([]byte("abc go gopher"))
	for _, r := range rs {
		fmt.Println(r)
		fmt.Printf("%#U\n", r)
	}
}

//十六进制 U+0061 'a' 6*16+1=97
