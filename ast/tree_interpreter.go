package ast

import (
	"reflect";
	"github.com/whatsthatsnail/simple_interpreter/lexer";
	"github.com/whatsthatsnail/simple_interpreter/errors"
)


type Interpreter struct{}

func (i Interpreter) visitBinary(b Binary) interface{} {
	left := b.X.Accept(i)
	right := b.Y.Accept(i)
	if reflect.TypeOf(left) == reflect.TypeOf(right) {
		switch left.(type) {
		case int:
			switch b.Op.TType {
			case lexer.PLUS:
				return left.(int) + right.(int)
			case lexer.MINUS:
				return left.(int) - right.(int)
			case lexer.STAR:
				return left.(int) * right.(int)
			case lexer.SLASH:
				return left.(int) / right.(int)
			}
		case float64:
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
		}
	} else {
		if b.Op.TType == lexer.PLUS {
			return left.(string) + right.(string)
		} else {
			errors.ThrowError(b.Op.Line, "Invalid string operation.")
			return nil
		}		
	}
	return nil
}

func (i Interpreter) visitUnary(u Unary) interface{} {
	switch u.X.Accept(i).(type) {
	case int:
		if u.Op.TType == lexer.MINUS {
			return 0 - u.X.Accept(i).(int)
		} else {
			return u.X.Accept(i).(int)
		}
	case float64:
		if u.Op.TType == lexer.MINUS {
			return 0 - u.X.Accept(i).(float64)
		} else {
			return u.X.Accept(i).(float64)
		}
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
