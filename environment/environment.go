package environment

import (
	"fmt";
	"github.com/whatsthatsnail/simple_interpreter/lexer";
	"github.com/whatsthatsnail/simple_interpreter/errors"
)

type Environment struct {
	Values map[string]interface{}
}

func NewEnvironment() Environment {
	env := Environment{}
	env.Values = make(map[string]interface{})
	return env
}

func (e *Environment) Get(name lexer.Token) interface{} {
	value, ok := e.Values[name.Lexeme]
	if ok {
		return value
	}

	// TODO: Make this a runtime error
	errors.ThrowError(name.Line, fmt.Sprintf("Undefined variable '%s'.", name.Lexeme))
	return nil
}

func (e *Environment) Assign(name lexer.Token, value interface{}) {
	_, ok := e.Values[name.Lexeme]
	if ok {
		e.Values[name.Lexeme] = value
	} else {
		// TODO: Make this a runtime error
		errors.ThrowError(name.Line, fmt.Sprintf("Undefined variable '%s'.", name.Lexeme))
	}
}