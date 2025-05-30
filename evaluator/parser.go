package evaluator

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Parser represents an expression parser
type Parser struct {
	tokens []string
	pos    int
}

// NewParser creates a new parser instance
func NewParser() *Parser {
	return &Parser{
		tokens: make([]string, 0),
		pos:    0,
	}
}

// Parse parses an expression string into an expression tree
func (p *Parser) Parse(expression string) (ExprNode, error) {
	// Tokenize the expression
	p.tokens = tokenize(expression)
	p.pos = 0

	// Parse the expression
	return p.parseExpression()
}

// parseExpression parses an expression: term (('+' | '-') term)*
func (p *Parser) parseExpression() (ExprNode, error) {
	expr, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	for p.pos < len(p.tokens) {
		op := p.tokens[p.pos]
		if op != "+" && op != "-" {
			break
		}
		p.pos++

		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}

		expr = &BinaryOpNode{
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}

	return expr, nil
}

// parseTerm parses a term: factor (('*' | '/') factor)*
func (p *Parser) parseTerm() (ExprNode, error) {
	expr, err := p.parseFactor()
	if err != nil {
		return nil, err
	}

	for p.pos < len(p.tokens) {
		op := p.tokens[p.pos]
		if op != "*" && op != "/" {
			break
		}
		p.pos++

		right, err := p.parseFactor()
		if err != nil {
			return nil, err
		}

		expr = &BinaryOpNode{
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}

	return expr, nil
}

// parseFactor parses a factor: number | '(' expression ')'
func (p *Parser) parseFactor() (ExprNode, error) {
	if p.pos >= len(p.tokens) {
		return nil, fmt.Errorf("unexpected end of expression")
	}

	token := p.tokens[p.pos]
	p.pos++

	if token == "(" {
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		if p.pos >= len(p.tokens) || p.tokens[p.pos] != ")" {
			return nil, fmt.Errorf("expected ')'")
		}
		p.pos++

		return expr, nil
	}

	// Try to parse as number
	value, err := strconv.ParseFloat(token, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid number: %s", token)
	}

	return &ValueNode{Value: value}, nil
}

// tokenize splits the expression into tokens
func tokenize(expression string) []string {
	var tokens []string
	var current strings.Builder

	for i := 0; i < len(expression); i++ {
		c := rune(expression[i])

		switch {
		case unicode.IsSpace(c):
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		case c == '(' || c == ')' || c == '+' || c == '-' || c == '*' || c == '/':
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			tokens = append(tokens, string(c))
		default:
			current.WriteRune(c)
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}
