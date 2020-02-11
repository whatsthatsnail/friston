package ast

import (
	"time"
)

// TODO: ClockFunc seems to not be a Function, fix this.

// Returns Unix time in seconds.
type ClockFunc struct{}

func (c ClockFunc) Arity() int { return 0 }

func (c ClockFunc) Call(i Interpreter, args interface{}) interface{} {
	now := time.Now()
	return float64(now.UnixNano()) / 1000000000
}
