package main

import (
	"bytes"
	"fmt"
)

//func (b *Buffer) ReadBytes(delim byte) (line []byte,err error)

func main() {
	var d byte = 'e' //分隔符  比如 可以用空格做分隔符
	buf := bytes.NewBufferString("你好esmieth")
	fmt.Println(buf.String()) //你好esmieth
	b, _ := buf.ReadBytes(d)  //读到分隔符，并返回给b
	fmt.Println(buf.String()) //smieth
	fmt.Println(string(b))    //你好e
}
