package main

import (
	"log"
	"os"
)

func main() {
	file, _ := os.Create("file.log")
	logger := log.New(file, "", log.Ldate|log.Ltime)
	logger.SetPrefix("[Info]")
	logger.Println("A debug message here")
	logger.SetPrefix("[Info]")
	logger.Println("An Info Message here")
	logger.SetFlags(logger.Flags() | log.LstdFlags)
	logger.Println("A different prefix")
	logger.SetPrefix("[Debug]")
	logger.Println("Here is debug")
}
