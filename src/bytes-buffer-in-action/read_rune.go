package main

import (
	"bytes"
	"fmt"
)

// func (b *Buffer) ReadRune() (r rune,size int,err error)

func main() {
	buf := bytes.NewBufferString("你好smith")
	fmt.Println(buf.String())
	b, n, _ := buf.ReadRune() //取出第一个rune
	fmt.Println(buf.String()) //好smith
	fmt.Println(string(b))    //你
	fmt.Println(n)            //3,"你“作为utf8存储占3个byte

	b, n, _ = buf.ReadRune()  //再取出一个rune
	fmt.Println(buf.String()) //smith
	fmt.Println(string(b))    //好
	fmt.Println(n)            //3
}
