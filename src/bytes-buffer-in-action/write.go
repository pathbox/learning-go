package main

import (
	"bytes"
	"fmt"
)

//func (b *Buffer) Write(p []byte) (n int,err error)
// Write方法，将一个byte类型的slice放到缓冲器的尾部
// 数据方向: slice => buffer slice
func main() {
	s := []byte(" world")
	buf := bytes.NewBufferString("hello")
	fmt.Println(buf.String()) //hello
	buf.Write(s)              //将s这个slice添加到buf的尾部
	fmt.Println(buf.String()) //hello world

}
