package main

import (
	"bytes"
	"fmt"
)

func main() {
	fmt.Printf("%s", bytes.TrimSpace([]byte(" \t\n a lone gopher \n\t\r\n")))
}
