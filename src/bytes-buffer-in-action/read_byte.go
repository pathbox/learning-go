package main

import (
	"bytes"
	"fmt"
)

//func (b *Buffer) ReadByte() (c byte,err error)
func main() {
	buf := bytes.NewBufferString("hello")
	fmt.Println(buf.String())
	b, _ := buf.ReadByte()    //取出第一个byte，赋值给b
	fmt.Println(buf.String()) //ello
	fmt.Println(string(b))    //h
}
