package ast

import (
	"reflect"
	"fmt"
	"friston/lexer"
	"friston/errors"
	env "friston/environment"
)

type Interpreter struct{
	Repl bool
	globals env.Environment
	environment env.Environment
}

func NewInterpreter(repl bool) Interpreter {
	i := Interpreter{}
	i.Repl = repl
	// Define global scope envionment (parent = nil)
	i.globals = env.NewEnvironment()
	
	i.globals.Declare("clock",  ClockFunc{})
	
	i.environment = i.globals
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

func (i Interpreter) visitLogic(l Logic) interface{} {
	left := i.evaluate(l.X)

	if l.Op.TType == lexer.OR {
		if isTruth(left) {
			return true
		} else {
			return isTruth(i.evaluate(l.Y))
		}
	}

	if l.Op.TType == lexer.AND {
		right := i.evaluate(l.Y)
		return isTruth(left) && isTruth(right)
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

	i.environment.Assign(a.Name.Lexeme, value)
	return value
}

func (i Interpreter) visitCall(c Call) interface{} {
	// Callee should probably be an IDENTIFIER, but really it can be anything, almost.
	callee := i.evaluate(c.Callee)

	var arguments []interface{}
	for _, arg := range(c.Arguments) {
		arguments = append(arguments, i.evaluate(arg))
	}

	// Cast the callee to type callable.function, and call it if it is a callable type.
	function, ok := callee.(Function)
	if !ok {
		// TODO: Runtime errors!
		errors.ThrowError(c.Paren.Line, "Can only call functions.")
	}
	
	// Check function arity. (Number of arguments)
	if len(arguments) != function.Arity() {
		errors.ThrowError(c.Paren.Line, "Expected " + string(function.Arity()) + " but got " + string(len(arguments)) + " arguments.")
	}

	return function.Call(i, arguments)
}

// Statement visitor methods:

func (i Interpreter) visitExprStmt(e ExprStmt) interface {} {
	value := i.evaluate(e.Expr)

	if i.Repl {
		fmt.Printf("%v\n", value)
	}

	return nil
}

func (i Interpreter) visitIfStmt(stmt IfStmt) interface {} {
	if isTruth(i.evaluate(stmt.Condition)) {
		i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		i.execute(stmt.ElseBranch)
	}

	return nil
}

func (i Interpreter) visitWhileStmt(stmt WhileStmt) interface {} {
	for isTruth(i.evaluate(stmt.Condition)) {
		i.execute(stmt.LoopBranch)
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

	i.environment.Declare(d.Name.Lexeme, value)
	return nil
}

func (i Interpreter) visitBlock(b Block) interface{} {
	// Create a new environment, enclosed by the current scope.
	enclosing := i.environment
	enclosed := env.NewEnvironment()
	enclosed.AddParent(enclosing)
	
	// Set the current environment to the inner scope. 
	i.environment = enclosed

	for _, s := range(b.Stmts) {
		i.execute(s)
	}


	return nil
}