package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// ItemType is type of lexer items
type ItemType int

func (typ ItemType) String() string {
	strType := "UNDEFINED"

	switch typ {
	case IERR:
		strType = "IERR"
	case INUMBER:
		strType = "INUMBER"
	case ILPAR:
		strType = "ILPAR"
	case IRPAR:
		strType = "IRPAR"
	case IADD:
		strType = "IADD"
	case ISUB:
		strType = "ISUB"
	case IMUL:
		strType = "IMUL"
	case IDIV:
		strType = "IDIV"
	case EOF:
		strType = "EOF"
	}
	return strType
}

// Item types from lexer items
const (
	// EOF Item type denoting end of input
	EOF ItemType = -1

	// lexing error occured
	IERR ItemType = iota

	// positive integer
	INUMBER

	// left parenthesis
	ILPAR

	// right parenthesis
	IRPAR

	// plus + symbol
	IADD

	// minus - symbol
	ISUB

	// multiply * symbol
	IMUL

	// divide / symbol
	IDIV
)

// Tokens understood by lexer
const (
	numbers   = "0123456789"
	operators = "+-*/"
	white     = " \n\r\t"
	lpar      = "("
	rpar      = ")"
)

// Type for lexing state machine, function returns another function which represents state transition
type stateFn func(*Lexer) stateFn

// LexItem are items emitted by lexer
type LexItem struct {
	Typ ItemType
	Pos int
	Val string
}

// Debug stringify lex item
func (li LexItem) String() string {
	return fmt.Sprintf("Type: %s, Val: %q, Pos: %d", li.Typ, li.Val, li.Pos)
}

// Lexer class containing state of lexer
type Lexer struct {
	// input text
	text string

	// start index (of input text) of currently lexing token
	start int

	// current position (of input text) of lexer
	Pos int

	// width of last read rune
	width int

	// output channel of lexer
	items chan LexItem
}

func (l *Lexer) dumpState() {
	fmt.Printf("%#v\n", l)
}

// Move to next ASCII or UTF-8 character/rune
func (l *Lexer) next() rune {
	if l.Pos >= len(l.text) {
		l.width = 0
		return -1
	}

	r, w := utf8.DecodeRuneInString(l.text[l.Pos:])
	l.width = w
	l.Pos += w

	return r
}

// Go back one character/rune
func (l *Lexer) backup() {
	l.Pos -= l.width
	_, w := utf8.DecodeLastRuneInString(l.text[:l.Pos])
	l.width = w
}

// Check what rune is next in input
func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// Consume next rune if its one of 'runes' parameter, otherwise no effect
func (l *Lexer) consume(runes string) bool {
	if strings.ContainsRune(runes, l.next()) {
		return true
	}
	l.backup()
	return false
}

// Consume all following runes that are one of 'runes'
func (l *Lexer) consumeAll(runes string) {
	for l.consume(runes) {
	}
}

// Ignores all un-emitted characters (moves start to current pos)
func (l *Lexer) ignore() {
	l.start = l.Pos
	l.width = 0
}

// Emits lexer item to output channel
func (l *Lexer) emit(typ ItemType) {
	l.items <- LexItem{
		Typ: typ,
		Pos: l.Pos,
		Val: l.text[l.start:l.Pos],
	}
	l.start = l.Pos
}

// Helper func that emits error item IERR with message
func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- LexItem{
		Pos: l.Pos,
		Typ: IERR,
		Val: fmt.Sprintf(format, args...),
	}
	return nil
}

// Starting state of state machine, peeks forward and decides what lexing function should be used
func lexFn(l *Lexer) stateFn {
	r := l.peek()
	switch {
	case r == -1:
		l.emit(EOF)
		return nil

	case strings.ContainsRune(white, r):
		return lexWhite
	case strings.ContainsRune(operators, r):
		return lexOperator
	case strings.ContainsRune(numbers, r):
		return lexNumber
	case strings.ContainsRune(lpar, r):
		return lexLpar
	case strings.ContainsRune(rpar, r):
		return lexRpar
	default:
		return l.errorf("Invalid symbol: %q", r)
	}
}

// Lexes operators
func lexOperator(l *Lexer) stateFn {
	op := l.next()
	switch op {
	case '+':
		l.emit(IADD)
	case '-':
		l.emit(ISUB)
	case '*':
		l.emit(IMUL)
	case '/':
		l.emit(IDIV)
	default:
		return l.errorf("lexOperator: inValid operator: %q", op)
	}

	return lexFn
}

// Lexes left parenthesis
func lexLpar(l *Lexer) stateFn {
	l.consume(lpar)
	l.emit(ILPAR)
	return lexFn
}

// Lexes right parenthesis
func lexRpar(l *Lexer) stateFn {
	l.consume(rpar)
	l.emit(IRPAR)
	return lexFn
}

// lexes numbers
func lexNumber(l *Lexer) stateFn {
	l.consumeAll(numbers)
	l.emit(INUMBER)
	return lexFn
}

// Lexes whitespaces and thrashes them (no emitting)
func lexWhite(l *Lexer) stateFn {
	l.consumeAll(white)
	l.ignore()
	return lexFn
}

// Items method gets channel of lex items
func (l *Lexer) Items() chan LexItem {
	return l.items
}

// Run method will is the core part of this lexer
// It fires off the state machine, starting with lexFn and takes return function as next state
// once some function returns nil (e.g. EOF), it stops lexing and returns
func (l *Lexer) Run() {
	defer close(l.items)

	for fun := lexFn; fun != nil; {
		fun = fun(l)
	}

}

// Lex is constructor for lexer
func Lex(text string) *Lexer {
	return &Lexer{
		items: make(chan LexItem),
		text:  text,
	}
}
