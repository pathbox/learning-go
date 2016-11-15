package main

import (
	"os"
)

func test() error {
	f, err := os.Create("text.txt")
	if err != nil {
		return err
	}

	defer f.Close()
	f.WriteString("Hello World!\n")
	return nil
}

func main() {
	test()
}
