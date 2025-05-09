package evaluator

import (
	"fmt"
)

// Expr represents an expression that can be evaluated
type Expr interface {
	Evaluate() (float64, error)
}

// BinaryExpr represents a binary operation (e.g., 1 + 2)
type BinaryExpr struct {
	Left     Expr
	Operator string
	Right    Expr
}

// LiteralExpr represents a number literal
type LiteralExpr struct {
	Value float64
}

// Evaluate implements the Expr interface for BinaryExpr
func (b *BinaryExpr) Evaluate() (float64, error) {
	left, err := b.Left.Evaluate()
	if err != nil {
		return 0, err
	}

	right, err := b.Right.Evaluate()
	if err != nil {
		return 0, err
	}

	switch b.Operator {
	case "+":
		return left + right, nil
	case "-":
		return left - right, nil
	case "*":
		return left * right, nil
	case "/":
		if right == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return left / right, nil
	default:
		return 0, fmt.Errorf("unknown operator: %s", b.Operator)
	}
}

// Evaluate implements the Expr interface for LiteralExpr
func (l *LiteralExpr) Evaluate() (float64, error) {
	return l.Value, nil
}
