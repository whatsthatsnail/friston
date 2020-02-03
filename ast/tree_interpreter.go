package ast

import (
	"reflect"
	"fmt";
	"github.com/whatsthatsnail/simple_interpreter/lexer";
	"github.com/whatsthatsnail/simple_interpreter/errors";
	env "github.com/whatsthatsnail/simple_interpreter/environment"
)

type Interpreter struct{
	Repl bool
	environment env.Environment
}

func NewInterpreter(repl bool) Interpreter {
	i := Interpreter{}
	i.Repl = repl
	// Define global scope envionment (parent = nil)
	i.environment = env.NewEnvironment()
	return i
}

func (i Interpreter) Interpret(stmts []Statement) {
	for _, s := range(stmts) {
		i.execute(s)
	}
}

// Helper methods:

func (i Interpreter) evaluate(expr Expression) interface{} {
	return expr.Accept(i)
}

func (i Interpreter) execute(stmt Statement) {
	stmt.Accept(i)
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

// Node visitor methods:

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
	case lexer.BANG_EQUAL:
		return !isEqual(left, right)
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

func (i Interpreter) visitVariable(vr Variable) interface{} {
	return i.environment.Get(vr.Name)
}

func (i Interpreter) visitAssignment(a Assignment) interface{} {
	value := i.evaluate(a.Value)

	i.environment.Assign(a.Name, value)
	return value
}

// Statement visitor methods:

func (i Interpreter) visitExprStmt(e ExprStmt) interface {} {
	value := i.evaluate(e.Expr)

	if i.Repl {
		fmt.Printf("%v\n", value)
	}

	return nil
}

func (i Interpreter) visitPrintStmt(p PrintStmt) interface {} {
	value := i.evaluate(p.Expr)
	fmt.Printf("%v\n", value)
	return nil
}

func (i Interpreter) visitVarDecl(d VarDecl) interface {} {
	var value interface{}
	if d.Initializer != nil {
		value = i.evaluate(d.Initializer)
	}

	i.environment.Values[d.Name.Lexeme] = value
	return nil
}

func (i Interpreter) visitBlock(b Block) interface{} {
	i.environment.AddParent(env.NewEnvironment())

	for _, s := range(b.Stmts) {
		i.execute(s)
	}

	return nil
}