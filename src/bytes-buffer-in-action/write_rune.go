package main

import (
	"bytes"
	"fmt"
)

// WriteRune方法，将一个rune类型的数据放到缓冲器的尾部
// func (b *Buffer) WriteRune(r Rune) (n int,err error)

func main() {
	var s rune = '好'
	buf := bytes.NewBufferString("hello")
	fmt.Println(buf.String()) //hello
	buf.WriteRune(s)
	fmt.Println(buf.String()) //hello好
}
