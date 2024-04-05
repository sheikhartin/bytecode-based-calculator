package parser

import (
	"fmt"
)

type TokenKind int

const (
	// Special tokens
	EOF = iota

	// Literals
	NUM
	IDENT

	// Operators
	FACT
	ADD
	SUB
	MUL
	DIV
	MOD
	POW

	// Symbols
	LPAREN
	RPAREN
	COMMA
	EQUAL
)

var tokenNames = map[TokenKind]string{
	EOF:    "EOF",
	NUM:    "NUM",
	IDENT:  "IDENT",
	FACT:   "FACT",
	ADD:    "ADD",
	SUB:    "SUB",
	MUL:    "MUL",
	DIV:    "DIV",
	MOD:    "MOD",
	POW:    "POW",
	LPAREN: "LPAREN",
	RPAREN: "RPAREN",
	COMMA:  "COMMA",
	EQUAL:  "EQUAL",
}

func (t TokenKind) String() string {
	return tokenNames[t]
}

type Position struct {
	Row, Col int
}

type Token struct {
	Pos   Position
	Kind  TokenKind
	Value string
}

func (t Token) String() string {
	return fmt.Sprintf(
		"%d:%d\t%s\t%s",
		t.Pos.Row,
		t.Pos.Col,
		t.Kind,
		t.Value,
	)
}
