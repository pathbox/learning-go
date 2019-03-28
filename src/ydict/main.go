package main

import (
	"fmt"
	"os"
)

var (
	proxy string
)

func main() {
	//Check & load .env file
	loadEnv()

	if len(os.Args) == 1 {
		displayUsage()
		os.Exit(0)
	}

	if len(os.Args) == 2 && os.Args[1] == "-h" {
		displayUsage()
		os.Exit(0)
	}

	words, withVoice, withMore, isQuiet := parseArgs(os.Args[1:]) // os.Args[0] 是命令，os.Args[1:] 之后的所有是 word
	fmt.Println(words, withVoice, withMore, isQuiet)
	query(words, withVoice, withMore, isQuiet, len(words) > 1)
}
