package evaluator

import (
	"fmt"
)

// Expr represents an expression that can be evaluated
type ExprNode interface {
	Evaluate() (float64, error)
}

// BinaryExpr represents a binary operation (e.g., 1 + 2)
type BinaryOpNode struct {
	Left     ExprNode
	Operator string
	Right    ExprNode
}

// LiteralExpr represents a number literal
type ValueNode struct {
	Value float64
}

// Evaluate implements the Expr interface for BinaryExpr
func (b *BinaryOpNode) Evaluate() (float64, error) {
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
func (l *ValueNode) Evaluate() (float64, error) {
	return l.Value, nil
}
