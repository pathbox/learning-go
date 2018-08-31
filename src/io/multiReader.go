package main

import (
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	r1 := strings.NewReader("first reader \n")
	r2 := strings.NewReader("second reader \n")
	r3 := strings.NewReader("third reader\n")
	r := io.MultiReader(r2, r1, r3) // 按顺序处理Reader，前面的先处理

	if _, err := io.Copy(os.Stdout, r); err != nil {
		log.Fatal(err)
	}

}
