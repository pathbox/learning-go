package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Mode())
		fmt.Println(file.Name())
		// fmt.Println(file.Sys())
		fmt.Println(file.Size())
		fmt.Println(file.ModTime())
	}
}
