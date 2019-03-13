package main

import "fmt"

type AstNodeType int

const (
	ASTNODE_LEAF AstNodeType = iota
	ASTNODE_ADD
	ASTNODE_SUB
	ASTNODE_MUL
	ASTNODE_DIV
)

func (t AstNodeType) String() string {
	switch t {
	case ASTNODE_LEAF:
		return "ASTNODE_LEAF"
	case ASTNODE_ADD:
		return "ASTNODE_ADD"
	case ASTNODE_SUB:
		return "ASTNODE_SUB"
	case ASTNODE_MUL:
		return "ASTNODE_MUL"
	case ASTNODE_DIV:
		return "ASTNODE_DIV"
	default:
		return "Unknown"
	}
}

type AstNode struct {
	Typ   AstNodeType
	Value *string

	Left  *AstNode
	Right *AstNode
}

func (node AstNode) String() string {
	printableVal := ""

	if node.Value != nil {
		printableVal = *node.Value
	}

	return fmt.Sprintf("AstNode: %s -> %s", node.Typ, printableVal)
}

func NewAstNode(typ AstNodeType, value *string) *AstNode {
	return &AstNode{
		Typ:   typ,
		Value: value,
	}
}
