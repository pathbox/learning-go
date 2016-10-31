package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.Create("writeAt.txt")
	if err != nil {
		panic(err)
	}

	defer file.Close()
	file.WriteString("Hello World!")
	n, err := file.WriteAt([]byte("good morning"), 12)
	if err != nil {
		panic(err)
	}
	fmt.Println(n)
}
