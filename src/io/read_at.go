package main

import (
	"fmt"
	"strings"
)

func main() {
	reader := strings.NewReader("Hello World")

	p := make([]byte, 6)

	n, err := reader.ReadAt(p, 2) // 从索引2开始读取6个字符
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s, %d\n", p, n)
}
