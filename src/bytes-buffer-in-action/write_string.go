package main

import (
	"bytes"
	"fmt"
)

// WriteString方法，把一个字符串放到缓冲器的尾部
//func (b *Buffer) WriteString(s string)(n int,err error)

func main() {
	s := " world"
	buf := bytes.NewBufferString("hello")
	fmt.Println(buf.String()) //hello
	buf.WriteString(s)        //将string写入到buf的尾部
	fmt.Println(buf.String()) //hello world
}
