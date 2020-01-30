package ast

import (
	"reflect"
	"fmt";
	"github.com/whatsthatsnail/simple_interpreter/lexer";
	"github.com/whatsthatsnail/simple_interpreter/errors"
)

// TODO: Handle bools, equalities, comparisons, etc

type Interpreter struct{}

func (i Interpreter) visitBinary(b Binary) interface{} {
	// Evaluate each side all the way down the tree
	left := b.X.Accept(i)
	right := b.Y.Accept(i)

	// When each side's type matches, do the correct operation
	if reflect.TypeOf(left) == reflect.TypeOf(right) {
		switch left.(type) {
		// Ints can convert to floats, but not the other way around, so convert everything to a float for simplicity
		case float64, int:
			switch b.Op.TType {
			case lexer.PLUS:
				return left.(float64) + right.(float64)
			case lexer.MINUS:
				return left.(float64) - right.(float64)
			case lexer.STAR:
				return left.(float64) * right.(float64)
			case lexer.SLASH:
				return left.(float64) / right.(float64)
			}

		// Concatenate strings only when they are added
		case string:
			if b.Op.TType == lexer.PLUS {
				return left.(string) + right.(string)
			} else {
				errors.ThrowError(b.Op.Line, "Invalid string operation.")
				return nil
			}
		}

	// If each side's type doesn't match, attempt to concatenate them as strings
	} else {
		if b.Op.TType == lexer.PLUS {
			return fmt.Sprintf("%v%v", left, right)
		} else {
			errors.ThrowError(b.Op.Line, "Invalid string operation.")
			return nil
		}		
	}
	return nil
}

func (i Interpreter) visitUnary(u Unary) interface{} {
	// Negate ints or floats if the operator matches '-'
	switch u.X.Accept(i).(type) {
	case int, float64:
		if u.Op.TType == lexer.MINUS {
			return 0 - u.X.Accept(i).(float64)
		} else {
			return u.X.Accept(i).(float64)
		}
	// Otherwise, continue down the tree.
	default:
		return u.X.Accept(i)
	}
}

func (i Interpreter) visitGroup(g Group) interface{} {
	return g.X.Accept(i)
}

func (i Interpreter) visitLiteral(l Literal) interface{} {
	return l.X.Literal
}
