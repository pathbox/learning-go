package main

import (
	"bytes"
	"fmt"
)

func main() {
	var ch byte
	fmt.Scanf("%c\n", &ch)

	buffer := new(bytes.Buffer)
	err := buffer.WriteByte(ch)
	if err == nil {
		fmt.Println("写入一个字节成功！准备读取该字节……")
		newCh, _ := buffer.ReadByte()
		fmt.Printf("读取的字节：%c\n", newCh)
	} else {
		fmt.Println("写入错误")
	}
}
