package ast

import (
	//"github.com/whatsthatsnail/simple_interpreter/lexer";
	"fmt"
)

type Expr interface {
	children() []Expr
}

type Binary struct {
	x Expr
	//op lexer.Token
	y Expr
} 

type Literal struct {
	token int
}

func (b Binary) children() []Expr {
	return append(b.x.children(), b.y.children()...)
}

func (l Literal) children() []Expr {
	return []Expr{l}
}

// (4 + 2) + 6
func Test() {
	lit1 := Literal{4}
	lit2 := Literal{2}

	bin1 := Binary{lit1, lit2}

	lit3 := Literal{6}
	bin2 := Binary{bin1, lit3}

	fmt.Println(bin1.children())
	fmt.Println(bin2.children())
}