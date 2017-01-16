package main

import (
	"log"
	"os"
)

func main() {
	file, err := os.OpenFile(
		"files/empty_new.txt",
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	byteSlice := []byte("Hello World!\n")
	bytesWritten, err := file.Write(byteSlice) // bytesWriteen is 写入的bytes的数量，这是覆盖的写入，会覆盖之前文件中的内容
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Wrote %d bytes.\n", bytesWritten)
}
