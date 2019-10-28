package main

import (
	"bytes"
	"fmt"
)

// WriteByte方法，将一个byte类型的数据放到缓冲器的尾部
//func (b *Buffer) WriteByte(c byte) error

func main() {
	var s byte = '?'
	buf := bytes.NewBufferString("hello")
	fmt.Println(buf.String()) //把buf的内容转换为string，hello
	buf.WriteByte(s)          //将s写到buf的尾部
	fmt.Println(buf.String()) //hello？
}
