package parser

import (
	"reflect"
	"testing"
)

func TestGenerateBytecodeForNumbers(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"358", []string{"PUSH_NUM\t358"}},
		{"2.7182", []string{"PUSH_NUM\t2.7182"}},
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
		g := NewBytecodeGenerator(p.Nodes)
		if !reflect.DeepEqual(g.Bytecode, tt.want) {
			t.Errorf("Failed to generate bytecode. Got `%v`, expected `%v`.", g.Bytecode, tt.want)
		}
	}
}

func TestGenerateBytecodeForCalls(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"func()", []string{"CALL_FUNC\tfunc\t0"}},
		{"func(22)", []string{"PUSH_NUM\t22", "CALL_FUNC\tfunc\t1"}},
		{"func(x, y,)", []string{"LOAD_VAR\tx", "LOAD_VAR\ty", "CALL_FUNC\tfunc\t2"}},
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
		g := NewBytecodeGenerator(p.Nodes)
		if !reflect.DeepEqual(g.Bytecode, tt.want) {
			t.Errorf("Failed to generate bytecode. Got `%v`, expected `%v`.", g.Bytecode, tt.want)
		}
	}
}

func TestGenerateBytecodeForUnaryOperations(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"! 6", []string{"PUSH_NUM\t6", "UNARY_OP\tFACT"}},
		{"! ! 3", []string{"PUSH_NUM\t3", "UNARY_OP\tFACT", "UNARY_OP\tFACT"}},
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
		g := NewBytecodeGenerator(p.Nodes)
		if !reflect.DeepEqual(g.Bytecode, tt.want) {
			t.Errorf("Failed to generate bytecode. Got `%v`, expected `%v`.", g.Bytecode, tt.want)
		}
	}
}

func TestGenerateBytecodeForBinaryOperations(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"% 25 5", []string{"PUSH_NUM\t25", "PUSH_NUM\t5", "BINARY_OP\tMOD"}},
		{"* 1 ! 0", []string{"PUSH_NUM\t1", "PUSH_NUM\t0", "UNARY_OP\tFACT", "BINARY_OP\tMUL"}},
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
		g := NewBytecodeGenerator(p.Nodes)
		if !reflect.DeepEqual(g.Bytecode, tt.want) {
			t.Errorf("Failed to generate bytecode. Got `%v`, expected `%v`.", g.Bytecode, tt.want)
		}
	}
}
