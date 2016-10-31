package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := strings.NewReader("Hello World!")
	reader.Seek(-6, os.SEEK_END)
	r, _, _ := reader.ReadRune()
	fmt.Println(r)
}
