package parser

import (
	"fmt"
	"strings"
)

type Parser struct {
	tokens  []Token
	currTok Token
	nextTok Token
	Nodes   []ASTNode
}

func (p Parser) String() string {
	var nodeStrings []string
	for _, node := range p.Nodes {
		if stringer, ok := node.(fmt.Stringer); ok {
			nodeStrings = append(nodeStrings, stringer.String())
		} else {
			nodeStrings = append(
				nodeStrings,
				fmt.Sprintf("(Node with no `String` method: %T)", node),
			)
		}
	}
	return strings.Join(nodeStrings, "; ")
}

func (p *Parser) advance() {
	if len(p.tokens) > 0 && p.currTok.Kind == 0 && p.nextTok.Kind == 0 {
		p.nextTok = p.tokens[0]
		p.tokens = p.tokens[1:]
	}
	p.currTok, p.nextTok = p.nextTok, Token{}
	if len(p.tokens) > 0 {
		p.nextTok = p.tokens[0]
		p.tokens = p.tokens[1:]
	}
}

func (p *Parser) expectKind(expected TokenKind) error {
	if p.currTok.Kind != expected {
		return fmt.Errorf("Expected `%s`, got `%s`!", expected, p.currTok.Kind)
	}
	p.advance()
	return nil
}

func (p *Parser) parseNumber() (*NumberNode, error) {
	num := &NumberNode{Value: p.currTok.Value}
	p.advance()
	return num, nil
}

func (p *Parser) parseIdentifier() *IdentifierNode {
	ident := &IdentifierNode{Value: p.currTok.Value}
	p.advance()
	return ident
}

func (p *Parser) parseAtomic() (ExprNode, error) {
	if p.currTok.Kind == NUM {
		return p.parseNumber()
	} else if p.currTok.Kind == IDENT {
		return p.parseIdentifier(), nil
	}
	return nil, nil
}

func (p *Parser) parseGroupedExpression() (ExprNode, error) {
	if err := p.expectKind(LPAREN); err != nil {
		return nil, err
	}
	node, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	if err := p.expectKind(RPAREN); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *Parser) parseCall() (*CallNode, error) {
	callee := p.parseIdentifier()
	p.expectKind(LPAREN)
	var args []ExprNode
	for {
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		} else if expr != nil {
			args = append(args, expr)
		}
		if p.currTok.Kind == COMMA {
			p.advance()
			if p.currTok.Kind == COMMA {
				return nil, fmt.Errorf(
					"Expected expression after comma at line %d, column %d.",
					p.currTok.Pos.Row,
					p.currTok.Pos.Col,
				)
			}
			continue
		} else {
			break
		}
	}
	p.expectKind(RPAREN)
	return &CallNode{Callee: callee, Args: args}, nil
}

func isUnaryOperator(kind TokenKind) bool {
	return kind == FACT
}

func (p *Parser) parseUnaryOperation() (ExprNode, error) {
	op := p.currTok.Kind
	p.advance()
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	return &UnaryOpNode{Operand: expr, Op: op}, nil
}

func isBinaryOperator(kind TokenKind) bool {
	return kind == ADD || kind == SUB || kind == MUL ||
		kind == DIV || kind == MOD || kind == POW
}

func (p *Parser) parseBinaryOperation() (ExprNode, error) {
	op := p.currTok.Kind
	p.advance()
	left, err := p.parseExpression()
	if err != nil {
		return nil, err
	} else if left == nil {
		return nil, fmt.Errorf(
			"Expected a left-hand operand for the `%s` operator at line %d, column %d.",
			op,
			p.currTok.Pos.Row,
			p.currTok.Pos.Col,
		)
	}
	right, err := p.parseExpression()
	if err != nil {
		return nil, err
	} else if right == nil {
		return nil, fmt.Errorf(
			"Expected a right-hand operand for the `%s` operator at line %d, column %d.",
			op,
			p.currTok.Pos.Row,
			p.currTok.Pos.Col,
		)
	}
	return &BinaryOpNode{Left: left, Op: op, Right: right}, nil
}

func (p *Parser) parseTerm() (ExprNode, error) {
	if p.currTok.Kind == LPAREN {
		return p.parseGroupedExpression()
	} else if p.currTok.Kind == IDENT && p.nextTok.Kind == LPAREN {
		return p.parseCall()
	} else if isUnaryOperator(p.currTok.Kind) {
		return p.parseUnaryOperation()
	}
	return p.parseAtomic()
}

func (p *Parser) parseExpression() (ExprNode, error) {
	if isBinaryOperator(p.currTok.Kind) {
		return p.parseBinaryOperation()
	}
	return p.parseTerm()
}

func (p *Parser) parseFullExpression() (ExprNode, error) {
	value, err := p.parseExpression()
	if err != nil {
		return nil, err
	} else if p.currTok.Kind != EOF {
		return nil, fmt.Errorf(
			"Unexpected token `%s` found after expression at line %d, column %d.",
			p.currTok.Value,
			p.currTok.Pos.Row,
			p.currTok.Pos.Col,
		)
	}
	return value, nil
}

func (p *Parser) parseVariableDeclaration() (*VariableDeclNode, error) {
	variable := p.parseIdentifier()
	if err := p.expectKind(EQUAL); err != nil {
		return nil, err
	}
	value, err := p.parseFullExpression()
	if err != nil {
		return nil, err
	}
	return &VariableDeclNode{Variable: variable, Value: value}, nil
}

func (p *Parser) parseStatement() (StmtNode, error) {
	if p.currTok.Kind == IDENT && p.nextTok.Kind == EQUAL {
		return p.parseVariableDeclaration()
	}
	return nil, nil
}

func (p *Parser) Parse() error {
	p.advance()
	stmt, err := p.parseStatement()
	if err != nil {
		return err
	} else if stmt != nil {
		p.Nodes = append(p.Nodes, stmt)
	}

	expr, err := p.parseFullExpression()
	if err != nil {
		return err
	} else if expr != nil {
		p.Nodes = append(p.Nodes, expr)
	}

	if stmt == nil && expr == nil {
		return fmt.Errorf(
			"Invalid grammar in position %d and %d.",
			p.currTok.Pos.Row,
			p.currTok.Pos.Col,
		)
	}
	return nil
}

func NewParser(tokens []Token) (*Parser, error) {
	p := &Parser{tokens: tokens}
	return p, p.Parse()
}
