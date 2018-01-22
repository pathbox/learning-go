package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("filetoread.txt")

	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize) // buffer 的大小就是文件的大小

	bytesize, err := file.Read(buffer)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("bytes read: ", bytesize)
	fmt.Println("bytestream to string: ", string(buffer))
}
