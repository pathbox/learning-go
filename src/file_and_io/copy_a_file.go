package main

import (
	"io"
	"log"
	"os"
)

// os.open os.create and io.copy
func main() {
	originalFile, err := os.Open("files/test.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer originalFile.Close()

	newFile, err := os.Create("files/test_copy.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer newFile.Close()
	// Copy the bytes to destination from source
	bytesWritten, err := io.Copy(newFile, originalFile)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Copied %d bytes.", bytesWritten)
	// Commit the file contents
	// Flushes memory to dis
	err = newFile.Sync()
	if err != nil {
		log.Fatal(err)
	}
}
