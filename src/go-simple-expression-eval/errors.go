package main

import "fmt"

// EvalErrorType error code type
type EvalErrorType int

const (
	// ErrParser error code raised by parser
	ErrParser EvalErrorType = iota

	// ErrLexer error code raised by lexer
	ErrLexer

	// ErrInterpreter error code raised by interpreter
	ErrInterpreter
)

func (err EvalErrorType) String() string {
	switch err {
	case ErrParser:
		return "Parser error"
	case ErrLexer:
		return "Lexer error"
	case ErrInterpreter:
		return "Interpreter error"
	default:
		return "UNDEFINED"
	}
}

// EvalError represents error during lexing, parsing or interpreting
type EvalError struct {
	// message
	s string

	// type of error
	code EvalErrorType
}

func (err EvalError) String() string {
	return fmt.Sprintf("%s: %s", err.code, err.s)
}

// NewLexerError intantiates new Lexer error
func NewLexerError(msg string, args ...interface{}) *EvalError {
	return &EvalError{
		s:    fmt.Sprintf(msg, args...),
		code: ErrLexer,
	}
}

// NewParserError intantiates new Parser error
func NewParserError(msg string, args ...interface{}) *EvalError {
	return &EvalError{
		s:    fmt.Sprintf(msg, args...),
		code: ErrParser,
	}
}

// NewInterpreterError intantiates new Interpreter error
func NewInterpreterError(msg string, args ...interface{}) *EvalError {
	return &EvalError{
		s:    fmt.Sprintf(msg, args...),
		code: ErrInterpreter,
	}
}
