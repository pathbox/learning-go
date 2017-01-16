package main

import (
	"log"
	"os"
)

func main() {
	file, err := os.Open("files/test.txt")
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	file, err = os.OpenFile("files/test.txt", os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
}

// Use these attributes individually or combined
// with an OR for second arg of OpenFile()
// e.g. os.O_CREATE|os.O_APPEND
// or os.O_CREATE|os.O_TRUNC|os.O_WRONLY

// os.O_RDONLY // Read only
// os.O_WRONLY // Write only
// os.O_RDWR // Read and write
// os.O_APPEND // Append to end of file
// os.O_CREATE // Create is none exist
// os.O_TRUNC // Truncate file when opening
