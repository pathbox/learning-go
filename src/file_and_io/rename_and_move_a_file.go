package main

import (
	"log"
	"os"
)

func main() {
	originalPath := "./files/empty.txt"
	newPath := "./files/empty_new.txt"

	err := os.Rename(originalPath, newPath)
	if err != nil {
		log.Fatalln(err)
	}
}
