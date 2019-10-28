package main

import (
	"bytes"
	"os"
)

// WriteTo方法，将一个缓冲器的数据写到w里，w是实现io.Writer的，比如os.File
// 数据方向:  buffer byte slice => io.Writer 数据从buffer中取出

func main() {
	file, _ := os.Create("text.txt")
	buf := bytes.NewBufferString("hello world")
	buf.WriteTo(file)
	//或者使用写入，fmt.Fprintf(file,buf.String())
}
