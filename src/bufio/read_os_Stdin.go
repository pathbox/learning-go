package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	for {
		line, ioErr := in.ReadString('\n')
		if ioErr != nil {
			fmt.Println("It is over")
			os.Exit(-1)
		}
		fmt.Println("echo: ", line)
	}

}
