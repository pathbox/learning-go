package main

import (
	"fmt"
	"go/scanner"
	"go/token"
)

func main() {
	scanner := scanner.Scanner{}
	source := []byte("n := 1\nfmt.Println(n)")
	errorHandler := func(_ token.Position, msg string) {
		fmt.Printf("error handler called: %s\n", msg)
	}
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(source))
	scanner.Init(file, source, errorHandler, 0)

	for {
		position, tok, literal := scanner.Scan()
		fmt.Printf("%d: %s", position, tok)
		if literal != ""{
			fmt.Printf(" %q", literal)
		}
		fmt.Println()
		if tok == token.EOF {
			break
		}
	}
}

