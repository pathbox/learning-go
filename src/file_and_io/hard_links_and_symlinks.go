package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	err := os.Link("files/test.txt", "test.txt")
	if err != nil {
		log.Println(err)
	}
	fmt.Println("creating sym")
	err = os.Symlink("files/test.txt", "test_sym.txt")
	if err != nil {
		log.Fatalln(err)
	}
	// Lstat will return file info, but if it is actually
	// a symlink, it will return info about the symlink.
	// It will not follow the link and give information
	// about the real file
	// Symlinks do not work in Windows
	fileInfo, err := os.Lstat("test_sym.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Link info: %+v", fileInfo)

	err = os.Lchown("test_sym.txt", os.Getuid(), os.Getgid())
	if err != nil {
		log.Panicln(err)
	}
}
