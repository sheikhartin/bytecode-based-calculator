package parser

import (
	"fmt"
)

type ASTNode interface {
	GenerateBytecode(g *BytecodeGenerator)
}

type ExprNode interface {
	ASTNode
}

type NumberNode struct {
	Value string
}

func (n NumberNode) String() string {
	return fmt.Sprintf("NumberNode{Value: %s}", n.Value)
}

func (n NumberNode) GenerateBytecode(g *BytecodeGenerator) {
	g.Emit("PUSH_NUM", n.Value)
}

type IdentifierNode struct {
	Value string
}

func (n IdentifierNode) String() string {
	return fmt.Sprintf("IdentifierNode{Value: %s}", n.Value)
}

func (n IdentifierNode) GenerateBytecode(g *BytecodeGenerator) {
	g.Emit("LOAD_VAR", n.Value)
}

type CallNode struct {
	Callee *IdentifierNode
	Args   []ExprNode
}

func (n CallNode) String() string {
	return fmt.Sprintf("CallNode{Callee: %s, Args: %s}", n.Callee, n.Args)
}

func (n CallNode) GenerateBytecode(g *BytecodeGenerator) {
	for _, arg := range n.Args {
		arg.GenerateBytecode(g)
	}
	g.Emit("CALL_FUNC", n.Callee.Value, fmt.Sprintf("%d", len(n.Args)))
}

type UnaryOpNode struct {
	Operand ExprNode
	Op      TokenKind
}

func (n UnaryOpNode) String() string {
	return fmt.Sprintf("UnaryOpNode{Operand: %s, Op: %s}", n.Operand, n.Op)
}

func (n UnaryOpNode) GenerateBytecode(g *BytecodeGenerator) {
	n.Operand.GenerateBytecode(g)
	g.Emit("UNARY_OP", n.Op.String())
}

type BinaryOpNode struct {
	Left  ExprNode
	Op    TokenKind
	Right ExprNode
}

func (n BinaryOpNode) String() string {
	return fmt.Sprintf("BinaryOpNode{Left: %s, Op: %s, Right: %s}", n.Left, n.Op, n.Right)
}

func (n BinaryOpNode) GenerateBytecode(g *BytecodeGenerator) {
	n.Left.GenerateBytecode(g)
	n.Right.GenerateBytecode(g)
	g.Emit("BINARY_OP", n.Op.String())
}

type StmtNode interface {
	ASTNode
}

type VariableDeclNode struct {
	Variable *IdentifierNode
	Value    ExprNode
}

func (n VariableDeclNode) String() string {
	return fmt.Sprintf("VariableDeclNode{Variable: %s, Value: %s}", n.Variable, n.Value)
}

func (n VariableDeclNode) GenerateBytecode(g *BytecodeGenerator) {
	n.Value.GenerateBytecode(g)
	g.Emit("STORE_VAR", n.Variable.Value)
}
