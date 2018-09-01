package main

import (
	"bytes"
	"fmt"
	"io"

	"strings"
)

func main() {
	r := strings.NewReader("some io.Reader stream to be read\n") // origin reader

	var buf bytes.Buffer // writer

	tee := io.TeeReader(r, &buf)

	p := make([]byte, 256)

	tee.Read(p) // Read action
	// data from origin reader to p and writer

	fmt.Println("p Read content: ", string(p))

	fmt.Println("Writer Recive: ", buf.String())

}

// 数据会从origin reader 读取到 p 和 writer， TeeReader将reader和writer之间建立起了联系，便捷的将reader的数据传到writer
//可以在你读一个 Reader 的同时，将数据写入到一个 Writer 中
