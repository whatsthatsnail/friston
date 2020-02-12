package callables

import (
	"fmt"
	"time"
)

var Natives = map[string]Function{
	"clock" : clockFunc{},
	"println" : printlnFunc{},
	"print" : printFunc{},
}

// Returns Unix time in seconds.
type clockFunc struct{}

func (c clockFunc) Arity() int { return 0 }

func (c clockFunc) Call(interpreter interface{}, args []interface{}) interface{} {
	now := time.Now()
	return float64(now.UnixNano()) / 1000000
}

type printFunc struct{}

func (p printFunc) Arity() int { return 1 }

func (p printFunc) Call(interpreter interface{}, args []interface{}) interface{} {
	fmt.Print(args[0])
	return nil
}

type printlnFunc struct{}

func (p printlnFunc) Arity() int { return 1 }

func (p printlnFunc) Call(interpreter interface{}, args []interface{}) interface{} {
	fmt.Println(args[0])
	return nil
}