package parser

import (
	"fmt"
	"strings"
)

type Lexer struct {
	Input  []string
	pos    Position
	currCh byte
	nextCh byte
	Tokens []Token
}

func (l Lexer) String() string {
	var tokenStrings []string
	for _, token := range l.Tokens {
		tokenStrings = append(tokenStrings, token.String())
	}
	return strings.Join(tokenStrings, "\n")
}

func (l *Lexer) advance() {
	if l.pos.Row >= len(l.Input) {
		l.currCh = 0
		l.nextCh = 0
		return
	} else if l.pos.Col >= len(l.Input[l.pos.Row]) {
		l.pos.Row++
		l.pos.Col = 0
		l.advance()
		return
	}

	l.currCh = l.Input[l.pos.Row][l.pos.Col]
	l.pos.Col++
	if l.pos.Col >= len(l.Input[l.pos.Row]) {
		l.nextCh = 0
	} else {
		l.nextCh = l.Input[l.pos.Row][l.pos.Col]
	}
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isOperator(ch byte) bool {
	return ch == '+' || ch == '-' || ch == '*' || ch == '/' ||
		ch == '%' || ch == '^' || ch == '!'
}

func classifyOperator(ch byte) (TokenKind, error) {
	switch ch {
	case '!':
		return FACT, nil
	case '+':
		return ADD, nil
	case '-':
		return SUB, nil
	case '*':
		return MUL, nil
	case '/':
		return DIV, nil
	case '%':
		return MOD, nil
	case '^':
		return POW, nil
	}
	return 0, fmt.Errorf("Invalid operator: `%c`", ch)
}

func isSymbol(ch byte) bool {
	return ch == '(' || ch == ')' || ch == ',' || ch == '=' || ch == ';'
}

func classifySymbol(ch byte) (TokenKind, error) {
	switch ch {
	case '(':
		return LPAREN, nil
	case ')':
		return RPAREN, nil
	case ',':
		return COMMA, nil
	case '=':
		return EQUAL, nil
	}
	return 0, fmt.Errorf("Invalid symbol: `%c`", ch)
}

func (l *Lexer) lexNumber() error {
	pos := l.pos
	var num string

	if l.currCh == '-' || l.currCh == '+' {
		num += string(l.currCh)
		l.advance()
	}
	for ; isDigit(l.currCh) || (l.currCh == '.' && !strings.Contains(num, ".")); l.advance() {
		num += string(l.currCh)
	}
	if l.currCh != 0 && !isWhitespace(l.currCh) && !isOperator(l.currCh) && !isSymbol(l.currCh) {
		return fmt.Errorf(
			"Invalid sequence `%s%c`! Line %d, column %d.",
			num,
			l.currCh,
			pos.Row+1,
			pos.Col,
		)
	}

	l.Tokens = append(l.Tokens, Token{
		Pos:   Position{Row: pos.Row + 1, Col: pos.Col},
		Kind:  NUM,
		Value: num,
	})
	return nil
}

func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}

func (l *Lexer) lexIdentifier() {
	pos := l.pos
	var id string

	for ; isLetter(l.currCh) || isDigit(l.currCh); l.advance() {
		id += string(l.currCh)
		if l.nextCh == 0 {
			l.advance()
			break
		}
	}
	l.Tokens = append(l.Tokens, Token{
		Pos:   Position{Row: pos.Row + 1, Col: pos.Col},
		Kind:  IDENT,
		Value: id,
	})
}

func (l *Lexer) Lex() error {
	for l.advance(); l.currCh != 0; {
		if isWhitespace(l.currCh) {
			l.advance()
			continue
		} else if isDigit(l.currCh) || (l.currCh == '-' || l.currCh == '+') && isDigit(l.nextCh) {
			if err := l.lexNumber(); err != nil {
				return err
			}
		} else if isLetter(l.currCh) {
			l.lexIdentifier()
		} else if isOperator(l.currCh) {
			kind, err := classifyOperator(l.currCh)
			if err != nil {
				return err
			}
			l.Tokens = append(l.Tokens, Token{
				Pos:   Position{Row: l.pos.Row + 1, Col: l.pos.Col},
				Kind:  kind,
				Value: string(l.currCh),
			})
			l.advance()
		} else if isSymbol(l.currCh) {
			if l.currCh == ';' && l.nextCh == ';' {
				l.pos.Row++
				l.pos.Col = 0
				l.advance()
				continue
			} else if l.currCh == ';' {
				return fmt.Errorf(
					"Invalid character `%c`! Line %d, column %d.",
					l.currCh,
					l.pos.Row+1,
					l.pos.Col,
				)
			}
			kind, err := classifySymbol(l.currCh)
			if err != nil {
				return err
			}
			l.Tokens = append(l.Tokens, Token{
				Pos:   Position{Row: l.pos.Row + 1, Col: l.pos.Col},
				Kind:  kind,
				Value: string(l.currCh),
			})
			l.advance()
		} else {
			return fmt.Errorf(
				"Invalid character `%c`! Line %d, column %d.",
				l.currCh,
				l.pos.Row+1,
				l.pos.Col,
			)
		}
	}
	l.Tokens = append(l.Tokens, Token{
		Pos:  Position{Row: l.pos.Row + 1, Col: 1},
		Kind: EOF,
	})
	return nil
}

func NewLexer(input string) (*Lexer, error) {
	l := &Lexer{Input: strings.Split(input, "\n")}
	return l, l.Lex()
}
