package main

import (
	"bytes"
	"fmt"
)

func main() {
	s := [][]byte{[]byte("foo"), []byte("bar"), []byte("baz")} // 二维[]byte数组

	js := bytes.Join(s, []byte(", ")) // 连接join数组元素
	fmt.Println(js)
	fmt.Printf("%s", js)
}
