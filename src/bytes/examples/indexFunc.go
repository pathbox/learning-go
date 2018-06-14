package main

import (
	"bytes"
	"fmt"
	"unicode"
)

func main() {
	f := func(c rune) bool {
		return unicode.Is(unicode.Han, c) // 是否是汉字编码
	}

	fmt.Println(bytes.IndexFunc([]byte("Hello 世界"), f))
	fmt.Println(bytes.IndexFunc([]byte("Hello, world"), f))
}
