package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Number float64

type Node interface {
	Eval() (Number, bool)
}

type Binary struct {
	op    byte
	left  Node
	right Node
}

func (n *Binary) Init(op byte, left, right Node) Node {
	n.op = op
	n.left = left
	n.right = right
	return n
}

func (n Binary) Eval() (Number, bool) {
	left, ok := n.left.Eval()
	if !ok {
		return 0, false
	}
	right, ok := n.right.Eval()
	if !ok {
		return 0, false
	}

	switch n.op {
	case '+':
		return left + right, true
	case '-':
		return left - right, true
	case '*':
		return left * right, true
	case '/':
		if right == 0 {
			return 0, false
		}
		return left / right, true
	}
	return 0, false
}

func (n *Binary) String() string {
	return fmt.Sprintf("(%s %c %s)", n.left, n.op, n.right)
}

type Leaf struct {
	value Number
}

func (n *Leaf) Init(value Number) Node {
	n.value = value
	return n
}

func (n *Leaf) Eval() (Number, bool) {
	return n.value, true
}

func (n *Leaf) String() string {
	return fmt.Sprintf("%v", n.value) // %v = default format
}

/* ==== Lexer ==== */
type Lexer struct {
	data string
	pos  int
	Kind int
	Num  Number
	Oper byte
}

const (
	ERR  = iota // error
	NUM         // number
	LPAR        // left parenthesis
	RPAR        // right parenthesis
	OP          // operator
)

func (lexer *Lexer) Init(data string) *Lexer {
	lexer.data = data
	lexer.pos = 0
	return lexer
}

func (l *Lexer) Next() int {
	n := len(l.data)
	l.Kind = ERR
	if l.pos < n {
		switch char := l.data[l.pos]; char {
		case '+', '-', '*', '/':
			l.pos++
			l.Kind = OP
			l.Oper = char
		case '(':
			l.pos++
			l.Kind = LPAR
			l.Oper = char
		case ')':
			l.pos++
			l.Kind = RPAR
			l.Oper = char
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.':
			var value Number = 0
			var divisor Number = 1
			for ; l.pos < n && '0' <= l.data[l.pos] && l.data[l.pos] <= '9'; l.pos++ {
				value = value*10 + Number(l.data[l.pos]-'0')
			}
			if l.pos < n && l.data[l.pos] == '.' { // 处理小数点后的数
				l.pos++
				for ; l.pos < n && '0' <= l.data[l.pos] && l.data[l.pos] <= '9'; l.pos++ {
					value = value*10 + Number(l.data[l.pos]-'0') // 计算.后的整数值和整数值合在一起
					divisor *= 10                                // 计算是小数点多少位
				}
			}
			l.Kind = NUM
			l.Num = value / divisor
		}
	}
	return l.Kind
}

/* ==== Parser ==== */
type Parser struct {
	lexer      *Lexer
	precedence map[byte]int
}

func (p *Parser) Init(data string) *Parser {
	p.lexer = new(Lexer).Init(data)
	p.precedence = make(map[byte]int)
	p.lexer.Next()
	return p
}

func (p *Parser) AddOperator(op byte, precedence int) {
	p.precedence[op] = precedence
}

func (p *Parser) Parse() (Node, bool) {
	lhs, ok := p.parsePrimary()
	if !ok {
		return nil, false
	}
	// starting with 1 instead of 0, because
	// map[*]int returns 0 for non-existant items
	node, ok := p.parseOperators(lhs, 1)
	if !ok {
		return nil, false
	}
	return node, true
}

func (p *Parser) parsePrimary() (Node, bool) {
	switch p.lexer.Kind {
	case NUM:
		node := new(Leaf).Init(p.lexer.Num)
		p.lexer.Next()
		return node, true
	case LPAR:
		p.lexer.Next()
		node, ok := p.Parse()
		if !ok {
			return nil, false
		}
		if p.lexer.Kind == RPAR {
			p.lexer.Next()
		}
		return node, true
	}
	return nil, false
}

func (p *Parser) parseOperators(lhs Node, min_precedence int) (Node, bool) {
	var ok bool
	var rhs Node
	for p.lexer.Kind == OP && p.precedence[p.lexer.Oper] >= min_precedence {
		op := p.lexer.Oper
		p.lexer.Next()
		rhs, ok = p.parsePrimary()
		if !ok {
			return nil, false
		}
		for p.lexer.Kind == OP && p.precedence[p.lexer.Oper] > p.precedence[op] {
			op2 := p.lexer.Oper
			rhs, ok = p.parseOperators(rhs, p.precedence[op2])
			if !ok {
				return nil, false
			}
		}
		lhs = new(Binary).Init(op, lhs, rhs)
	}
	return lhs, true
}

func main() {
	var node Node
	var result Number
	var p *Parser
	var parseOk, evalOk bool
	in := bufio.NewReader(os.Stdin)
	line, ioErr := in.ReadString('\n')
	for len(line) > 0 {
		line = strings.TrimSpace(line)
		fmt.Printf("Read: %q\n", line) // %q = quoted string
		p = new(Parser).Init(line)
		p.AddOperator('+', 1)
		p.AddOperator('-', 1)
		p.AddOperator('*', 2)
		p.AddOperator('/', 2)
		node, parseOk = p.Parse()
		if parseOk {
			fmt.Printf("Parsed: %s\n", node)
			result, evalOk = node.Eval()
			if evalOk {
				fmt.Printf("Evaluated: %v\n", result) // %v = default format
			} else {
				fmt.Printf("%s = Evaluation error\n", line)
			}
		} else {
			fmt.Printf("%s = Syntax error\n", line)
		}
		if ioErr != nil {
			return
		}
		line, ioErr = in.ReadString('\n')
	}
}
