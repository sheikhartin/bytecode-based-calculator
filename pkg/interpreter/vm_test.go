package interpreter

import (
	"reflect"
	"testing"

	"github.com/sheikhartin/bytecode-based-calculator/pkg/parser"
)

func TestProcessNumbers(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"358", 358},
		{"2.7182", 2.7182},
	}

	for _, tt := range tests {
		l, err := parser.NewLexer(tt.input)
		if err != nil {
			t.Fatalf("Failed to tokenize input `%s`: %v", tt.input, err)
			continue
		}
		p, err := parser.NewParser(l.Tokens)
		if err != nil {
			t.Fatalf("Failed to initialize parser with tokens from input `%s`: %v", tt.input, err)
			continue
		}
		g := parser.NewBytecodeGenerator(p.Nodes)
		vm := NewVM()
		if err := vm.Execute(g.Bytecode); err != nil {
			t.Fatalf("Execution error for input `%s`: %v", tt.input, err)
		} else if got := vm.Stack[len(vm.Stack)-1]; !reflect.DeepEqual(got, tt.want) {
			t.Errorf(
				"The execution output does not match the expectations! Input `%s`, got `%v`, want `%v`.",
				tt.input,
				got,
				tt.want,
			)
		}
	}
}

func TestProcessFunctionCalls(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"fact(3)", 6},
		{"min(5, 10, 15, 20)", 5},
		{"max(5, 10, 15, 20)", 20},
	}

	for _, tt := range tests {
		l, err := parser.NewLexer(tt.input)
		if err != nil {
			t.Fatalf("Failed to tokenize input `%s`: %v", tt.input, err)
			continue
		}
		p, err := parser.NewParser(l.Tokens)
		if err != nil {
			t.Fatalf("Failed to initialize parser with tokens from input `%s`: %v", tt.input, err)
			continue
		}
		g := parser.NewBytecodeGenerator(p.Nodes)
		vm := NewVM()
		if err := vm.Execute(g.Bytecode); err != nil {
			t.Fatalf("Execution error for input `%s`: %v", tt.input, err)
		} else if got := vm.Stack[len(vm.Stack)-1]; !reflect.DeepEqual(got, tt.want) {
			t.Errorf(
				"The execution output does not match the expectations! Input `%s`, got `%v`, want `%v`.",
				tt.input,
				got,
				tt.want,
			)
		}
	}
}

func TestProcessUnaryOperations(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"! 5", 120},
		{"! ! 0", 1},
	}

	for _, tt := range tests {
		l, err := parser.NewLexer(tt.input)
		if err != nil {
			t.Fatalf("Failed to tokenize input `%s`: %v", tt.input, err)
			continue
		}
		p, err := parser.NewParser(l.Tokens)
		if err != nil {
			t.Fatalf("Failed to initialize parser with tokens from input `%s`: %v", tt.input, err)
			continue
		}
		g := parser.NewBytecodeGenerator(p.Nodes)
		vm := NewVM()
		if err := vm.Execute(g.Bytecode); err != nil {
			t.Fatalf("Execution error for input `%s`: %v", tt.input, err)
		} else if got := vm.Stack[len(vm.Stack)-1]; !reflect.DeepEqual(got, tt.want) {
			t.Errorf(
				"The execution output does not match the expectations! Input `%s`, got `%v`, want `%v`.",
				tt.input,
				got,
				tt.want,
			)
		}
	}
}

func TestProcessBinaryOperations(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"- PI E", 0.423310825130748},
		{"% - 50 10 2", 0},
	}

	for _, tt := range tests {
		l, err := parser.NewLexer(tt.input)
		if err != nil {
			t.Fatalf("Failed to tokenize input `%s`: %v", tt.input, err)
			continue
		}
		p, err := parser.NewParser(l.Tokens)
		if err != nil {
			t.Fatalf("Failed to initialize parser with tokens from input `%s`: %v", tt.input, err)
			continue
		}
		g := parser.NewBytecodeGenerator(p.Nodes)
		vm := NewVM()
		if err := vm.Execute(g.Bytecode); err != nil {
			t.Fatalf("Execution error for input `%s`: %v", tt.input, err)
		} else if got := vm.Stack[len(vm.Stack)-1]; !reflect.DeepEqual(got, tt.want) {
			t.Errorf(
				"The execution output does not match the expectations! Input `%s`, got `%v`, want `%v`.",
				tt.input,
				got,
				tt.want,
			)
		}
	}
}

func TestProcessVariableDeclaration(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"x = 10", 10},
		{"y = x", 10},
	}

	vm := NewVM() // To keep the variable in the next rounds
	for _, tt := range tests {
		l, err := parser.NewLexer(tt.input)
		if err != nil {
			t.Fatalf("Failed to tokenize input `%s`: %v", tt.input, err)
			continue
		}
		p, err := parser.NewParser(l.Tokens)
		if err != nil {
			t.Fatalf("Failed to initialize parser with tokens from input `%s`: %v", tt.input, err)
			continue
		}
		g := parser.NewBytecodeGenerator(p.Nodes)
		if err := vm.Execute(g.Bytecode); err != nil {
			t.Fatalf("Execution error for input `%s`: %v", tt.input, err)
		} else if got := vm.Stack[len(vm.Stack)-1]; !reflect.DeepEqual(got, tt.want) {
			t.Errorf(
				"The execution output does not match the expectations! Input `%s`, got `%v`, want `%v`.",
				tt.input,
				got,
				tt.want,
			)
		}
	}
}
