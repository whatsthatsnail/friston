package interpreter

import (
	"fmt"
	"time"
)

var Natives = map[string]Function{
	"clock" : clockNative{},
	"println" : printlnNative{},
	"print" : printNative{},
}

// Returns Unix time in seconds.
type clockNative struct{}

func (c clockNative) Arity() int { return 0 }

func (c clockNative) Call(i Interpreter, args []interface{}) interface{} {
	now := time.Now()
	return float64(now.UnixNano()) / 1000000000
}

type printNative struct{}

func (p printNative) Arity() int { return 1 }

func (p printNative) Call(i Interpreter, args []interface{}) interface{} {
	fmt.Print(args[0])
	return nil
}

type printlnNative struct{}

func (p printlnNative) Arity() int { return 1 }

func (p printlnNative) Call(i Interpreter, args []interface{}) interface{} {
	fmt.Println(args[0])
	return nil
}