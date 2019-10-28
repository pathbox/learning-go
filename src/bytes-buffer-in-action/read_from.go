package main

import (
	"bytes"
	"fmt"
	"os"
)

// ReadFrom方法，从一个实现io.Reader接口的r，把r的内容读到缓冲器里，n返回读的数量
func main() {
	file, _ := os.Open("text.txt")
	buf := bytes.NewBufferString("bob ")
	buf.ReadFrom(file)        // 从file中读取全部内容到buffer slice 中
	fmt.Println(buf.String()) //bob hello world
}
