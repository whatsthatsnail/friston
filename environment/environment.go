package environment

import (
	"fmt"
	"friston/errors"
	"friston/lexer"
)

type Environment struct {
	Values map[string]interface{}
	parent *Environment
}

func NewEnvironment() Environment {
	env := Environment{}
	env.Values = make(map[string]interface{})
	return env
}

// Creates a new environment with the specified parent environment.
func NewEnclosed(parent Environment) Environment {
	enclosed := NewEnvironment()
	enclosed.AddParent(parent)
	return enclosed
}

func (e *Environment) GetParent() interface{} {
	if e.parent != nil {
		return *e.parent
	}

	return "Global"
}

func (e *Environment) AddParent(parentEnv Environment) {
	e.parent = &parentEnv
}

func (e *Environment) Get(name lexer.Token) interface{} {
	value, ok := e.Values[name.Lexeme]
	if ok {
		return value
	}

	// Recursively check parents for variable if it's not found in current scope.
	if e.parent != nil {
		return e.parent.Get(name)
	}

	// TODO: Make this a runtime error
	errors.ThrowError(name.Line, fmt.Sprintf("Undefined variable '%s'.", name.Lexeme))
	return nil
}

// Assign value to a variable in current scope, or parent scopes, if it exists.
func (e *Environment) Assign(name lexer.Token, value interface{}) {
	_, ok := e.Values[name.Lexeme]
	if ok {
		e.Values[name.Lexeme] = value
		return
	}

	if e.parent != nil {
		// Just like Get(), recurively check parent scopes for the target variable.
		e.parent.Assign(name, value)
		return
	}

	// TODO: Stop execution with runtime error.
	// TODO: Fix line number.
	// TODO: Convert name back to token, not string, in all uses.
	errors.ThrowError(0, fmt.Sprintf("Undefined variable '%s'.", name.Lexeme))
}

// Declare a new variable in the current scope
func (e *Environment) Declare(name string, value interface{}) {
	e.Values[name] = value
}
