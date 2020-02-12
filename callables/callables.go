package callables

type Function interface {
	Call(interpreter interface{}, args []interface{}) interface{}
	Arity() int
}