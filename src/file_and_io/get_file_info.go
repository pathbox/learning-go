package main

import (
	"fmt"
	"log"
	"os"
)

var (
	file *os.FileInfo
	err  error
)

func main() {
	fileInfo, err := os.Stat("./files/test.txt")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("File name: ", fileInfo.Name())
	fmt.Println("Size in bytes: ", fileInfo.Size())
	fmt.Println("Permissions: ", fileInfo.Mode())
	fmt.Println("Last modified: ", fileInfo.ModTime())
	fmt.Println("Is Directory: ", fileInfo.IsDir())
	fmt.Println("System interface type: %T\n", fileInfo.Sys())
	fmt.Printf("System interface type: %+v\n\n", fileInfo.Sys())
}
