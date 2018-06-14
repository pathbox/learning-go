package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	var b bytes.Buffer
	b.Write([]byte("Hello ")) // 将数据存到 buffer 中
	fmt.Fprintf(&b, "world!")
	b.WriteTo(os.Stdout) // 将buffer中的数据 输出
}
