package main

import (
	"bytes"
	"fmt"
)

// func (b *Buffer) Read(p []byte)(n int,err error)
// 数据方向： buffer slice => other byte slice
func main() {
	s1 := []byte("hello")
	buff := bytes.NewBuffer(s1)
	s2 := []byte(" world")
	buff.Write(s2)             //buffer 尾部追加
	fmt.Println(buff.String()) //hello world

	s3 := make([]byte, 3)
	buff.Read(s3)              //把buff的内容读入到s3，s3的容量为3，读了3个过来
	fmt.Println(buff.String()) //lo world
	fmt.Println(string(s3))    //hel
	buff.Read(s3)              //继续读入3个，原来的被覆盖

	fmt.Println(buff.String()) //world
	fmt.Println(len(s3))
	fmt.Println(string(s3)) //"lo "
}
