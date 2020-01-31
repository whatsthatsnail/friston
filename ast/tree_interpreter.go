package ast

import (
	"reflect"
	"fmt";
	"github.com/whatsthatsnail/simple_interpreter/lexer";
	"github.com/whatsthatsnail/simple_interpreter/errors"
)

type Interpreter struct{}

// Helper methods:

func (i Interpreter) evaluate(expr Expression) interface{} {
	return expr.Accept(i)
}

// Nil, false bools, zero, empty strings are false, all else is true.
func isTruth(expr interface{}) bool {
	switch expr.(type) {
	case nil:
		return false
	case bool:
		if expr.(bool) == false {
			return false
		} else {
			return true
		}
	case int, float64:
		if expr.(float64) == 0 {
			return false
		} else {
			return true
		}
	case string:
		if len(expr.(string)) == 0 {
			return false
		} else {
			return true
		}
	default:
		return true
	}
}

// Nil is only equal to itself (our equality differs from Golang).
func isEqual(left interface{}, right interface{}) bool {
	if left == nil && right == nil {
		return true
	} else if left == nil {
		return false
	}

	return left == right
}

func checkNumberOperand(operator lexer.Token, number interface{}) bool {
	switch number.(type) {
	case int, float64:
		return true
	default:
		errors.ThrowError(operator.Line, "Operand must be a number.")
		return false
	}
}

func checkNumberOperands(operator lexer.Token, left interface{}, right interface{}) bool {
	if reflect.TypeOf(left) == reflect.TypeOf(right) {
		switch left.(type) {
		case int, float64:
			return true
		default:
			errors.ThrowError(operator.Line, "Operands must be a number.")
			return false
		}
	} else {
		errors.ThrowError(operator.Line, "Operand types must match.")
		return false
	}
}

// Visitor methods:

func (i Interpreter) visitBinary(b Binary) interface{} {
	// Evaluate each side all the way down the tree.
	left := b.X.Accept(i)
	right := b.Y.Accept(i)

	switch b.Op.TType {
	// Basic arithmetic:
	case lexer.MINUS:
		if checkNumberOperands(b.Op, left, right) {
			return left.(float64) - right.(float64)
		}
	case lexer.STAR:
		if checkNumberOperands(b.Op, left, right) {
			return left.(float64) * right.(float64)
		}
	case lexer.SLASH:
		if checkNumberOperands(b.Op, left, right) {
			if right.(float64) == 0 {
				errors.ThrowError(b.Op.Line, "Division by zero.")
				return nil
			}
			return left.(float64) / right.(float64)
		}

	// Addition (includes string concatenation):
	case lexer.PLUS:
		switch left.(type) {
		case int, float64:
			if checkNumberOperands(b.Op, left, right) {
				return left.(float64) + right.(float64)
			}
		case string:
			return fmt.Sprintf("%v%v", left, right)
		}

	// Comparisons:
	case lexer.GREATER:
		if checkNumberOperands(b.Op, left, right) {
			return left.(float64) > right.(float64)
		}
	case lexer.GREATER_EQUAL:
		if checkNumberOperands(b.Op, left, right) {
			return left.(float64) >= right.(float64)
		}
	case lexer.LESS:
		if checkNumberOperands(b.Op, left, right) {
			return left.(float64) < right.(float64)
		}
	case lexer.LESS_EQUAL:
		if checkNumberOperands(b.Op, left, right) {
			return left.(float64) <= right.(float64)
		}
	case lexer.EQUAL_EQUAL:
		return isEqual(left, right)
	}

	// Unreachable.
	return nil
}

func (i Interpreter) visitUnary(u Unary) interface{} {
	right := i.evaluate(u.X)

	switch u.Op.TType {
	case lexer.MINUS:
		if checkNumberOperand(u.Op, right) {
			return -right.(float64)
		}
	case lexer.BANG:
		return !isTruth(right)
	}

	// Unreachable.
	return nil
}

func (i Interpreter) visitGroup(g Group) interface{} {
	return i.evaluate(g.X)
}

func (i Interpreter) visitLiteral(l Literal) interface{} {
	return l.X.Literal
}
