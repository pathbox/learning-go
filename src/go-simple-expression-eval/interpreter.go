package main

import (
	"strconv"
)

// Implementation of ADD operation
func add(a int, b int) int {
	return a + b
}

// Implementation of SUBTRACT operation
func sub(a int, b int) int {
	return a - b
}

// Implementation of MULTIPLY operation
func mul(a int, b int) int {
	return a * b
}

// Implementation of DIVIDE operation
func div(a int, b int) int {
	return a / b
}

// Type of arithmetic function taking two ints and returning ints
type arithmeticFunc func(int, int) int

// type for mapping AST node types to corresponding arithmetic operation
type interpretMap map[AstNodeType]arithmeticFunc

// Recursive post order traversal that evaluates AST
// 1. Visit left
// 2. Visit right
// 3. Visit self
func postOrderTraversal(node *AstNode, functions interpretMap) (int, *EvalError) {
	if node == nil {
		return 0, NewInterpreterError("Expected evaluatable node, got nil")
	}

	// If we are on number node
	if node.Typ == ASTNODE_LEAF {
		if node.Value == nil {
			return 0, NewInterpreterError("Expected value, got nil")
		}

		// Parse string val to integer
		number, err := strconv.Atoi(*node.Value)
		if err != nil {
			return 0, NewInterpreterError("Unable to parse number %s", err)
		}

		// return it to higher stack frame (numbers should occur only in leaf nodes)
		return number, nil
	}

	aritFunc := functions[node.Typ]
	left, err := postOrderTraversal(node.Left, functions)
	if err != nil {
		return 0, err
	}

	right, err := postOrderTraversal(node.Right, functions)
	if err != nil {
		return 0, err
	}

	// use its value to do computation
	return aritFunc(left, right), nil
}

// Interpret is function that evaluates AST and returns corresponding result
func Interpret(ast *AstNode) (int, *EvalError) {
	var astInterpretMap = interpretMap{
		ASTNODE_ADD: add,
		ASTNODE_SUB: sub,
		ASTNODE_MUL: mul,
		ASTNODE_DIV: div,
	}

	return postOrderTraversal(ast, astInterpretMap)
}
