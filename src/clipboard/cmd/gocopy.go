package main

import (
	"fmt"

	"github.com/atotto/clipboard"
)

func main() {

	text, err := clipboard.ReadAll()
	fmt.Println("from: ", text)

	text = "Have a nice day"
	err = clipboard.WriteAll(text)
	if err != nil {
		panic(err)
	}

	fmt.Println("done")

}
