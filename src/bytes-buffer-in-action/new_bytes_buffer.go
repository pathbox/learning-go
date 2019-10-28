package main

import (
	"bytes"
	"fmt"
)

func main() {
	buf1 := bytes.NewBufferString("hello")
	buf2 := bytes.NewBuffer([]byte("hello"))
	buf3 := bytes.NewBuffer([]byte{'h', 'e', 'l', 'l', 'o'})
	buf4 := bytes.NewBufferString("")
	buf5 := bytes.NewBuffer([]byte{})

	// var buf bytes.Buffer  定义一个buffer类型  buf为buffer struct 拥有其接口方法
	fmt.Println(buf1.String(), buf2.String(), buf3.String(), buf4, buf5, 1)
}
