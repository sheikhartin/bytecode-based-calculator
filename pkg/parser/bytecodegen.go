package parser

import (
	"fmt"
	"strings"
)

type BytecodeGenerator struct {
	ast      []ASTNode
	Bytecode []string
}

func (g BytecodeGenerator) String() string {
	return strings.Join(g.Bytecode, "\n")
}

func (g *BytecodeGenerator) Emit(op string, operands ...string) {
	g.Bytecode = append(g.Bytecode, fmt.Sprintf("%s\t%s", op, strings.Join(operands, "\t")))
}

func (g *BytecodeGenerator) Generate() {
	for _, node := range g.ast {
		node.GenerateBytecode(g)
	}
}

func NewBytecodeGenerator(ast []ASTNode) *BytecodeGenerator {
	g := &BytecodeGenerator{ast: ast}
	g.Generate()
	return g
}
