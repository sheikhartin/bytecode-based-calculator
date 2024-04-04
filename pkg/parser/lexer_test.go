package parser

import (
	"reflect"
	"testing"
)

func TestNumbers(t *testing.T) {
	input := "42 3.14 -7 +8.90"
	want := []Token{
		{Pos: Position{Row: 1, Col: 1}, Kind: NUM, Value: "42"},
		{Pos: Position{Row: 1, Col: 4}, Kind: NUM, Value: "3.14"},
		{Pos: Position{Row: 1, Col: 9}, Kind: NUM, Value: "-7"},
		{Pos: Position{Row: 1, Col: 12}, Kind: NUM, Value: "+8.90"},
		{Pos: Position{Row: 2, Col: 1}, Kind: EOF},
	}

	l, err := NewLexer(input)
	if err != nil {
		t.Fatalf("An error while lexing! %v", err)
	}
	if !reflect.DeepEqual(l.Tokens, want) {
		t.Errorf("It did not meet expectations!")
	}
}

func TestIdentifiers(t *testing.T) {
	input := "x x2 xy xy2 x0y0z0"
	want := []Token{
		{Pos: Position{Row: 1, Col: 1}, Kind: IDENT, Value: "x"},
		{Pos: Position{Row: 1, Col: 3}, Kind: IDENT, Value: "x2"},
		{Pos: Position{Row: 1, Col: 6}, Kind: IDENT, Value: "xy"},
		{Pos: Position{Row: 1, Col: 9}, Kind: IDENT, Value: "xy2"},
		{Pos: Position{Row: 1, Col: 13}, Kind: IDENT, Value: "x0y0z0"},
		{Pos: Position{Row: 2, Col: 1}, Kind: EOF},
	}

	l, err := NewLexer(input)
	if err != nil {
		t.Fatalf("An error while lexing! %v", err)
	}
	if !reflect.DeepEqual(l.Tokens, want) {
		t.Errorf("It did not meet expectations!")
	}
}

func TestOperators(t *testing.T) {
	input := "! + - * / % ^"
	want := []Token{
		{Pos: Position{Row: 1, Col: 1}, Kind: FACT, Value: "!"},
		{Pos: Position{Row: 1, Col: 3}, Kind: ADD, Value: "+"},
		{Pos: Position{Row: 1, Col: 5}, Kind: SUB, Value: "-"},
		{Pos: Position{Row: 1, Col: 7}, Kind: MUL, Value: "*"},
		{Pos: Position{Row: 1, Col: 9}, Kind: DIV, Value: "/"},
		{Pos: Position{Row: 1, Col: 11}, Kind: MOD, Value: "%"},
		{Pos: Position{Row: 1, Col: 13}, Kind: POW, Value: "^"},
		{Pos: Position{Row: 2, Col: 1}, Kind: EOF},
	}

	l, err := NewLexer(input)
	if err != nil {
		t.Fatalf("An error while lexing! %v", err)
	}
	if !reflect.DeepEqual(l.Tokens, want) {
		t.Errorf("It did not meet expectations!")
	}
}

func TestSymbols(t *testing.T) {
	input := "( ) , = ;; A worthless comment!"
	want := []Token{
		{Pos: Position{Row: 1, Col: 1}, Kind: LPAREN, Value: "("},
		{Pos: Position{Row: 1, Col: 3}, Kind: RPAREN, Value: ")"},
		{Pos: Position{Row: 1, Col: 5}, Kind: COMMA, Value: ","},
		{Pos: Position{Row: 1, Col: 7}, Kind: EQUAL, Value: "="},
		{Pos: Position{Row: 2, Col: 1}, Kind: EOF},
	}

	l, err := NewLexer(input)
	if err != nil {
		t.Fatalf("An error while lexing! %v", err)
	}
	if !reflect.DeepEqual(l.Tokens, want) {
		t.Errorf("It did not meet expectations!")
	}
}

func TestIllegalCharacters(t *testing.T) {
	input := "π ≈ ᭠˙Ꞌꜭ"

	_, err := NewLexer(input)
	if err == nil {
		t.Errorf("An error while lexing! %v", err)
	}
}
