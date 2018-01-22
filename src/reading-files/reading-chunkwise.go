package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	const BufferSize = 100
	file, err := os.Open("filetoread.txt")

	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	buffer := make([]byte, BufferSize) // 以定义的BufferSize 作为缓冲大小

	for {
		bytesread, err := file.Read(buffer) // 把数据读到buffer

		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
		fmt.Println("bytes read: ", bytesread)
		fmt.Println("bytestream to string: ", string(buffer[:bytesread]))
		// fmt.Println("bytestream to string: ", string(buffer[:]))
	}

}
