package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strconv"
)

var tests = []string{
	"(1+3)*7", // 28, example from task description.
	"1+3*7",   // 22, shows operator precedence.
	"7",       // 7, a single literal is a valid expression.
	"7/3",     // eval only does integer math.
	"7.3",     // this parses, but we disallow it in eval.
	"7^3",     // parses, but disallowed in eval.
	"go",      // a valid keyword, not valid in an expression.
	"3@7",     // error message is "illegal character."
	"",        // EOF seems a reasonable error message.
}

func main() {
	for _, exp := range tests {
		if r, err := parseAndEval(exp); err == nil {
			fmt.Println(exp, "=", r)
		} else {
			fmt.Printf("%s: %v\n", exp, err)
		}
	}
}

func parseAndEval(exp string) (int, error) {
	tree, err := parser.ParseExpr(exp)
	if err != nil {
		return 0, err
	}
	return eval(tree)
}

func eval(tree ast.Expr) (int, error) {
	switch n := tree.(type) {
	case *ast.BasicLit:
		if n.Kind != token.INT {
			return unsup(n.Kind)
		}
		i, _ := strconv.Atoi(n.Value)
		return i, nil
	case *ast.BinaryExpr:
		switch n.Op {
		case token.ADD, token.SUB, token.MUL, token.QUO:
		default:
			return unsup(n.Op)
		}
		x, err := eval(n.X)
		if err != nil {
			return 0, err
		}
		y, err := eval(n.Y)
		if err != nil {
			return 0, err
		}
		switch n.Op {
		case token.ADD:
			return x + y, nil
		case token.SUB:
			return x - y, nil
		case token.MUL:
			return x * y, nil
		case token.QUO:
			return x / y, nil
		}
	case *ast.ParenExpr:
		return eval(n.X)
	}
	return unsup(reflect.TypeOf(tree))
}

func unsup(i interface{}) (int, error) {
	return 0, errors.New(fmt.Sprintf("%v unsupported", i))
}
