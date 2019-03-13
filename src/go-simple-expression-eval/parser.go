package main

import (
	"container/list"
	"fmt"
)

// Determines precedence of operator given
func precedence(typ ItemType) int {
	switch typ {
	case IADD:
		fallthrough
	case ISUB:
		return 1

	case IMUL:
		fallthrough
	case IDIV:
		return 2

	default:
		return -1
	}
}

// Converts infix output form of Lexer to postfix form
func toPostfix(lx *Lexer) (*list.List, *EvalError) {

	opStack := NewStack()
	postFix := list.New()

	for item := range lx.Items() {
		// end of tok stream
		if item.Typ == EOF {
			continue
		}

		// lexing error
		if item.Typ == IERR {
			return nil, NewLexerError("at %d: %s", item.Pos, item.Val)
		}

		// if its number put to output
		if item.Typ == INUMBER {
			postFix.PushBack(item)
			continue
		}

		// if left parenth put to output
		if item.Typ == ILPAR {
			opStack.Push(item)
			continue
		}

		// if right parenth
		if item.Typ == IRPAR {
			// pop stack to output until we find left parenth in stack
			for opStack.Len() > 0 && opStack.Top().(LexItem).Typ != ILPAR {
				postFix.PushBack(opStack.Pop())
			}

			// if there is none then there is error in parity
			if opStack.Len() > 0 && opStack.Top().(LexItem).Typ != ILPAR {
				return nil, NewParserError("at %d: Unmatched paretheses", opStack.Top().(LexItem).Pos)
			}

			// we are in rparenth so if there is no lparenh its parity error
			if opStack.Len() == 0 {
				return nil, NewParserError("at %d: Missing '('", item.Pos)
			}
			// otherwise just trash it
			opStack.Pop()

		} else {
			// is any other operator
			// check precedence
			for opStack.Len() > 0 && precedence(item.Typ) <= precedence(opStack.Top().(LexItem).Typ) {
				// just put it to output
				postFix.PushBack(opStack.Pop())
			}
			// put it to stack
			opStack.Push(item)
		}
	}

	// empty stack to output
	for opStack.Len() > 0 {
		postFix.PushBack(opStack.Pop())
	}

	return postFix, nil
}

// helper method that translates Lexer item types to AST node types
func translateLexToAstType(typ ItemType) (AstNodeType, *EvalError) {
	switch typ {
	case IADD:
		return ASTNODE_ADD, nil
	case ISUB:
		return ASTNODE_SUB, nil
	case IMUL:
		return ASTNODE_MUL, nil
	case IDIV:
		return ASTNODE_DIV, nil
	default:
		return 0, NewParserError("Unexpected item type occured during parsing %q", typ)
	}
}

// Takes list of postfix formed lexer items and builds binary expression tree
func constructAst(postfixList *list.List) (*AstNode, *EvalError) {
	// stack for storing nodes for later computation
	stack := NewStack()

	// go trough all items
	for item := postfixList.Front(); item != nil; item = item.Next() {
		lexItem := item.Value.(LexItem)
		// if its number, create node and push it to stack
		if lexItem.Typ == INUMBER {
			stack.Push(NewAstNode(ASTNODE_LEAF, &lexItem.Val))
		} else {
			// otherwise convert type
			nodeType, err := translateLexToAstType(lexItem.Typ)
			if err != nil {
				return nil, NewParserError("at %d: Missing ')'", lexItem.Pos)
			}
			// create new note
			node := NewAstNode(nodeType, nil)

			// validate we have at least two items in stack
			if stack.Len() < 2 {
				return nil, NewParserError("at %d: Missing operand", lexItem.Pos)
			}

			// order important, otherwise we switch operands
			// Pop first time to Right operand
			node.Right = stack.Pop().(*AstNode)
			// Pop second time to Left operand
			node.Left = stack.Pop().(*AstNode)

			// push new node to stack
			stack.Push(node)
		}
	}

	// might occur when user inputs "()" expression, no root node
	if stack.Len() < 1 {
		return nil, NewParserError("Expression without root")
	}

	// pop last item from stack, its the root node of AST
	return stack.Pop().(*AstNode), nil
}

// helper debug func
func traversePreorder(root *AstNode) {
	if root == nil {
		return
	}

	fmt.Println(root)
	traversePreorder(root.Left)
	traversePreorder(root.Right)
}

// helper debug func
func traverseInorder(root *AstNode) {
	if root == nil {
		return
	}

	traversePreorder(root.Left)
	fmt.Println(root)
	traversePreorder(root.Right)
}

// helper debug func
func traversePostorder(root *AstNode) {
	if root == nil {
		return
	}

	traversePreorder(root.Left)
	traversePreorder(root.Right)
	fmt.Println(root)
}

// Parse parses given infix expression and produces Abstract syntax tree
func Parse(expr string) (*AstNode, *EvalError) {
	lx := Lex(expr)
	go lx.Run()

	postfixNotation, err := toPostfix(lx)

	if err != nil {
		return nil, err
	}

	abstractSyntaxTree, err := constructAst(postfixNotation)

	if err != nil {
		return nil, err
	}

	return abstractSyntaxTree, nil
}
