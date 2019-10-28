package main

import (
	"bytes"
	"fmt"
)

// Next方法，返回前n个byte（slice），原缓冲器变小
//func (b *Buffer) Next(n int) []byte
func main() {
	buf := bytes.NewBufferString("hello world")
	fmt.Println(buf.String())
	b := buf.Next(2)          //取前2个
	fmt.Println(buf.String()) //llo world
	fmt.Println(string(b))    //he
	buf.Next(2)
	buf.Next(2)
	buf.Next(2)
	fmt.Println(buf.String())
}
