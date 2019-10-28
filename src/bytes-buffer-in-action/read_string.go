package main

import (
	"bytes"
	"fmt"
)

//func (b *Buffer) ReadString(delim byte) (line string,err error)

func main() {
	var d byte = 'e'
	buf := bytes.NewBufferString("你好esmieth")
	fmt.Println(buf.String()) //你好esmieth
	b, _ := buf.ReadString(d) //读取到分隔符，并返回给b
	fmt.Println(buf.String()) //smieth
	fmt.Println(string(b))    //你好e
}
