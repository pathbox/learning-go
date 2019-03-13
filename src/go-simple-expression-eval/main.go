package main

import (
	"fmt"
	"os"
)

func main() {

	args := os.Args

	if len(args) < 2 {
		fmt.Println("Specify expressions to evaluate...\ne.g.: 1+2*(6-8)")
		return
	}

	// parse given expression into AST
	ast, err := Parse(args[1])

	if err != nil {
		fmt.Println(err)
		return
	}

	// Interpret AST
	result, err := Interpret(ast)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(result)
}
