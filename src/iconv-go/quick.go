package main

import (
	"fmt"
	iconv "github.com/djimenez/iconv-go"
)

func main() {
	// output, _ := iconv.ConvertString("Hello World!", fromEncoding, toEncoding)
	output, _ := iconv.ConvertString("Hello World!", "utf-8", "windows-1252")
	fmt.Println(output)

	converter, _ := iconv.NewConverter("utf-8", "gbk")
	outputput, _ := converter.ConvertString("Hello World!")
	fmt.Println(outputput)

	converter.Close()

	// in := []byte("Hello World!")
	// out := make([]byte, len(input))

	// bytesRead, bytesWritten, err := iconv.Convert(in, out, "utf-8", "latin1")
	// in := []byte("Hello World!")
	// out := make([]byte, len(input))

	// bytesRead, bytesWritten, err := iconv.Convert(in, out, "utf-8", "latin1")
}
