package visitors

import (
	"fmt"
	"reflect"
	"friston/ast"
	"friston/callables"
	"friston/environment"
	"friston/errors"
	"friston/lexer"
)

type Interpreter struct{
	Repl bool
	globals environment.Environment
	environment environment.Environment
}

func NewInterpreter(repl bool) Interpreter {
	i := Interpreter{}
	i.Repl = repl
	// Define global scope envionment (parent = nil)
	i.globals = environment.NewEnvironment()

	// Declare all native functions in the global environment
	for k, v := range callables.Natives {
		i.globals.Declare(k, v)
	}
	
	i.environment = i.globals
	return i
}

func (i Interpreter) Interpret(stmts []ast.Statement) {
	for _, s := range(stmts) {
		i.execute(s)
	}
}

// Helper methods:

func (i Interpreter) evaluate(expr ast.Expression) interface{} {
	return expr.Accept(i)
}

func (i Interpreter) execute(stmt ast.Statement) {
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
		}
	case int, float64:
		if expr.(float64) == 0 {
			return false
		}
	case string:
		if len(expr.(string)) == 0 {
			return false
		}
	default:
		return true
	}

	// Unreachable.
	return true
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

// Node Visitor methods:

func (i Interpreter) VisitBinary(b ast.Binary) interface{} {
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

func (i Interpreter) VisitLogic(l ast.Logic) interface{} {
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

func (i Interpreter) VisitUnary(u ast.Unary) interface{} {
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

func (i Interpreter) VisitGroup(g ast.Group) interface{} {
	return i.evaluate(g.X)
}

func (i Interpreter) VisitLiteral(l ast.Literal) interface{} {
	return l.X.Literal
}

func (i Interpreter) VisitVariable(vr ast.Variable) interface{} {
	return i.environment.Get(vr.Name)
}

func (i Interpreter) VisitAssignment(a ast.Assignment) interface{} {
	value := i.evaluate(a.Value)

	i.environment.Assign(a.Name, value)
	return value
}

func (i Interpreter) VisitCall(c ast.Call) interface{} {
	// Callee should probably be an IDENTIFIER, but really it can be anything, almost.
	callee := i.evaluate(c.Callee)

	var arguments []interface{}
	for _, arg := range(c.Arguments) {
		arguments = append(arguments, i.evaluate(arg))
	}

	// Cast the callee to type callable.function, and call it if it is a callable type.
	function, ok := callee.(callables.Function)
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

// Statement Visitor methods:

func (i Interpreter) VisitExprStmt(e ast.ExprStmt) {
	value := i.evaluate(e.Expr)

	if i.Repl {
		fmt.Printf("%v\n", value)
	}
}

func (i Interpreter) VisitIfStmt(stmt ast.IfStmt) {
	if isTruth(i.evaluate(stmt.Condition)) {
		i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		i.execute(stmt.ElseBranch)
	}
}

func (i Interpreter) VisitWhileStmt(stmt ast.WhileStmt) {
	for isTruth(i.evaluate(stmt.Condition)) {
		i.execute(stmt.LoopBranch)
	}
}

func (i Interpreter) VisitVarDecl(d ast.VarDecl) {
	var value interface{}
	if d.Initializer != nil {
		value = i.evaluate(d.Initializer)
	}

	i.environment.Declare(d.Name.Lexeme, value)
}

func (i Interpreter) VisitBlock(b ast.Block) {
	// Create a new environment, enclosed by the current scope.
	enclosing := i.environment
	enclosed := environment.NewEnvironment()
	enclosed.AddParent(enclosing)
	
	// Set the current environment to the inner scope. 
	i.environment = enclosed

	for _, s := range(b.Stmts) {
		i.execute(s)
	}
}