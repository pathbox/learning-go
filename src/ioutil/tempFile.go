package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	content := []byte("temporary file's content")
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(tmpfile.Name())

	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}

	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}
}
