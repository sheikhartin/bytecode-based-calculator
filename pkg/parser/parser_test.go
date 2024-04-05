package parser

import (
	"reflect"
	"testing"
)

func TestParseNumbers(t *testing.T) {
	tests := []struct {
		input string
		want  []ASTNode
	}{
		{"42", []ASTNode{&NumberNode{Value: "42"}}},
		{"3.14", []ASTNode{&NumberNode{Value: "3.14"}}},
	}

	for _, tt := range tests {
		l, err := NewLexer(tt.input)
		if err != nil {
			t.Fatalf("Failed to tokenize input `%s`: %v", tt.input, err)
			continue
		}
		p, err := NewParser(l.Tokens)
		if err != nil {
			t.Fatalf("Failed to initialize parser with tokens from input `%s`: %v", tt.input, err)
			continue
		}
		if !reflect.DeepEqual(p.Nodes, tt.want) {
			t.Errorf("Failed to parse expression. Got `%v`, expected `%v`.", p.Nodes, tt.want)
		}
	}
}

func TestParseUnaryOperations(t *testing.T) {
	tests := []struct {
		input string
		want  []ASTNode
	}{
		{"! 8", []ASTNode{&UnaryOpNode{Operand: &NumberNode{Value: "8"}, Op: FACT}}},
		{"! ! 4", []ASTNode{
			&UnaryOpNode{
				Operand: &UnaryOpNode{Operand: &NumberNode{Value: "4"}, Op: FACT},
				Op:      FACT,
			},
		}},
	}

	for _, tt := range tests {
		l, err := NewLexer(tt.input)
		if err != nil {
			t.Fatalf("Failed to tokenize input `%s`: %v", tt.input, err)
		}
		p, err := NewParser(l.Tokens)
		if err != nil {
			t.Fatalf("Failed to initialize parser with tokens from input `%s`: %v", tt.input, err)
		}
		if !reflect.DeepEqual(p.Nodes, tt.want) {
			t.Errorf("Failed to parse binary operation. Got `%v`, expected `%v`.", p.Nodes, tt.want)
		}
	}
}

func TestParseBinaryOperations(t *testing.T) {
	tests := []struct {
		input string
		want  []ASTNode
	}{
		{"+ 30 55.0", []ASTNode{
			&BinaryOpNode{
				Left:  &NumberNode{Value: "30"},
				Op:    ADD,
				Right: &NumberNode{Value: "55.0"},
			},
		}},
		{"^ 2 ! 3", []ASTNode{
			&BinaryOpNode{
				Left:  &NumberNode{Value: "2"},
				Op:    POW,
				Right: &UnaryOpNode{Operand: &NumberNode{Value: "3"}, Op: FACT},
			},
		}},
		{"^ 1000 ! 0", []ASTNode{
			&BinaryOpNode{
				Left:  &NumberNode{Value: "1000"},
				Op:    POW,
				Right: &UnaryOpNode{Operand: &NumberNode{Value: "0"}, Op: FACT},
			},
		}},
	}

	for _, tt := range tests {
		l, err := NewLexer(tt.input)
		if err != nil {
			t.Fatalf("Failed to tokenize input `%s`: %v", tt.input, err)
		}
		p, err := NewParser(l.Tokens)
		if err != nil {
			t.Fatalf("Failed to initialize parser with tokens from input `%s`: %v", tt.input, err)
		}
		if !reflect.DeepEqual(p.Nodes, tt.want) {
			t.Errorf("Failed to parse binary operation. Got `%v`, expected `%v`.", p.Nodes, tt.want)
		}
	}
}

func TestParseVariableDeclaration(t *testing.T) {
	tests := []struct {
		input string
		want  []ASTNode
	}{
		{"pi = 3.14", []ASTNode{
			&VariableDeclNode{
				Variable: &IdentifierNode{Value: "pi"},
				Value:    &NumberNode{Value: "3.14"},
			},
		}},
		{"currYear = - ^ 2 11 24", []ASTNode{
			&VariableDeclNode{
				Variable: &IdentifierNode{Value: "currYear"},
				Value: &BinaryOpNode{
					Left: &BinaryOpNode{
						Left:  &NumberNode{Value: "2"},
						Op:    POW,
						Right: &NumberNode{Value: "11"},
					},
					Op:    SUB,
					Right: &NumberNode{Value: "24"},
				},
			},
		}},
	}

	for _, tt := range tests {
		l, err := NewLexer(tt.input)
		if err != nil {
			t.Fatalf("Failed to tokenize input `%s`: %v", tt.input, err)
		}
		p, err := NewParser(l.Tokens)
		if err != nil {
			t.Fatalf("Failed to initialize parser with tokens from input `%s`: %v", tt.input, err)
		}
		if !reflect.DeepEqual(p.Nodes, tt.want) {
			t.Errorf("Failed to parse variable declaration. Got `%v`, expected `%v`.", p.Nodes, tt.want)
		}
	}
}
