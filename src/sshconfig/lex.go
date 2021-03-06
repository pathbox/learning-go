package sshconfig

import (
	"fmt"
	"strings"
	"unicode"
)

type pos int

type item struct {
	typ itemType
	pos pos
	val string
}

func (i item) String() string {
	switch {
	case i.typ == itemEOF:
		return "EOF"
	case i.typ == itemError:
		return i.val
	case len(i.val) > 10:
		return fmt.Sprintf("%.10q...", i.val)
	}
	return fmt.Sprintf("%q", i.val)
}

type itemType int

const (
	itemError itemType = iota
	itemEOF
	itemValue
	itemHost
	itemHostValue
	itemHostName
	itemUser
	itemPort
	itemProxyCommand
	itemHostKeyAlgorithms
	itemIdentityFile
)

// variables
var variables = map[string]itemType{
	"host":              itemHost,
	"hostname":          itemHostName,
	"user":              itemUser,
	"port":              itemPort,
	"proxycommand":      itemProxyCommand,
	"hostkeyalgorithms": itemHostKeyAlgorithms,
	"identityfile":      itemIdentityFile,
}

const eof = -1

// stateFn represents the state of the scanner as a function that returns the next state
type stateFn func(*lexer) stateFn

type lexer struct {
	input   string
	state   stateFn
	pos     pos
	start   pos
	width   pos
	lastPos pos
	items   chan item // channel of scanned items
}

func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := uft8.DecodeRuneInString(l.input[l.pos:])
	l.width = pos(w)
	l.pos += l.width
	return r
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) backup() {
	l.pos -= l.width
}

// emit passes an item back to the client
func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.start, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{itemError, l.start, fmt.Sprintf(format, args...)}
	return nil
}

func (l *lexer) nextItem() item {
	item := <-l.items
	l.lastPos = item.pos
	return item
}

func lex(input string) *lexer {
	l := &lexer{
		input: input,
		items: make(chan item),
	}
	go l.run()
	return l
}

func (l *lexer) run() {
	for l.state = lexEnv; l.state != nil; {
		l.state = l.state(l)
	}
}

func lexEnv(l *lexer) stateFn {
	switch r := l.next(); {
	case r == eof:
		l.emit(itemEOF)
		return nil
	case isAlphaNumeric(r):
		return lexVariable
	case r == '#':
		return lexComment
	case r == '\r':
		l.ignore()
		if l.next() != '\n' {
			l.errorf("expected \\n")
			return nil
		}
		l.ignore()
		return lexEnv
	case r == '\t' || r == ' ' || r == '\n':
		l.ignore()
		return lexEnv
	default:
		l.errorf("unable to parse character: %c", r)
		return nil
	}
}

func lexComment(l *lexer) stateFn {
	for {
		switch l.next() {
		case '\r':
			l.ignore()
			if l.next() != '\n' {
				l.errorf("expected \\n")
				return nil
			}
			l.ignore()
			return lexEnv
		case '\n', eof:
			l.ignore()
			return lexEnv
		default:
			l.ignore()
		}
	}
}

func lexVariable(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case isAlphaNumeric(r):
			// absorb
		case r == ' ' || r == '=':
			l.backup()
			variable := strings.ToLower(l.input[l.start:l.pos])

			if _, ok := variables[variable]; ok {
				l.emit(variables[variable])
				l.next()
				l.ignore()
				if variable == "host" {
					return lexHostValue
				}
				return lexValue
			}
			return lexValue
		default:
			pattern := l.input[l.start:l.pos]
			return l.errorf("invalid pattern: %#v", pattern)
		}
	}
}

func lexHostValue(l *lexer) stateFn {
	for {
		switch l.next() {
		case ' ':
			switch l.peek() {
			case '\r':
				if l.peek() != '\n' {
					return l.errorf("expected \\n")
				}
				fallthrough
			case '\n', eof:
				break
			default:
				// more coming, wait
				continue
			}
			l.emit(itemValue)
		case '\r':
			if l.peek() != '\n' {
				return l.errorf("expected \\n")
			}
			fallthrough
		case '\n':
			l.backup()
			l.emit(itemHostValue)
			return lexEnv
		case eof:
			l.backup()
			l.emit(itemHostValue)
			l.next()
			l.emit(itemEOF)
			return nil
		}
	}
}

func lexValue(l *lexer) stateFn {
	for {
		switch l.next() {
		case '\r':
			if l.peek() != '\n' {
				return l.errorf("expected \\n")
			}
			fallthrough
		case '\n':
			l.backup()
			l.emit(itemValue)
			return lexEnv
		case eof:
			l.backup()
			l.emit(itemValue)
			l.next()
			l.emit(itemEOF)
			return nil
		}
	}
}

// isAlphaNumeric reports whether r is an alphabetic or digit.
func isAlphaNumeric(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}
