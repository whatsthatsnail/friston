package environment

import (
	"fmt";
	"github.com/whatsthatsnail/simple_interpreter/lexer";
	"github.com/whatsthatsnail/simple_interpreter/errors"
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

func (e *Environment) Assign(name lexer.Token, value interface{}) {
	_, ok := e.Values[name.Lexeme]
	if ok {
		e.Values[name.Lexeme] = value
		return
	}

	// Just like Get(), recurively check parent scopes for the target variable.
	if e.parent != nil {
		e.parent.Assign(name, value)
		return
	}

	// TODO: Make this a runtime error
	errors.ThrowError(name.Line, fmt.Sprintf("Undefined variable '%s'.", name.Lexeme))
}
