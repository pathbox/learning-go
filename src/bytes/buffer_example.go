package main

import (
	"bytes"
	"fmt"
)

func main() {
	rd := bytes.NewBufferString("Hello World!")
	buf := make([]byte, 6)
	// 获取数据切片
	fmt.Println(len(buf))
	b := rd.Bytes()
	// 读出一部分数据，看看切片有没有变化
	fmt.Println(rd.String())
	fmt.Printf("%s\n", rd.String()) // World!
	fmt.Printf("%s\n\n", b)         // Hello World!
	fmt.Println("======", rd)
	// 写入一部分数据，看看切片有没有变化
	rd.Write([]byte("abcdefg"))
	fmt.Printf("%s\n", rd.String()) // World!abcdefg
	fmt.Printf("%s\n\n", b)         // Hello World!

	// 再读出一部分数据，看看切片有没有变化
	rd.Read(buf)
	fmt.Printf("%s\n", rd.String()) // abcdefg
	fmt.Printf("%s\n", b)           // Hello World!
	fmt.Println(rd)
}
