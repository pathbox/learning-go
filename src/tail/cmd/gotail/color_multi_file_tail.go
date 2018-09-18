package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gongo/9t"
)

func main() {
	filenames := os.Args[1:]
	fmt.Println(filenames)
	runner, err := ninetail.Runner(filenames, ninetail.Config{Colorize: true})
	if err != nil {
		log.Fatal(err)
	}
	runner.Run()
}
