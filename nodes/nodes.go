package nodes

import (
	"fmt"
	"log"
	"strings"
)

type Node interface {
	Children() []Node
	String() string
	Eval(ctx Context) interface{}
}

type Token interface {
	Value() interface{}
}

type Context interface {
	Get(s string) interface{}
}

type Internal func(Context, ...Node) interface{}

type Expr struct {
	children []Node
}

func (e *Expr) Children() []Node {
	return e.children
}

func (e *Expr) Eval(ctx Context) interface{} {
	log.Print("Expr:Eval", e)
	if len(e.children) <= 0 {
		return nil
	}

	fn := e.children[0].Eval(ctx)
	return fn.(Internal)(ctx, e.children[1:]...)
}

func (e *Expr) String() string {
	var elt []string
	for _, node := range e.children {
		elt = append(elt, node.String())
	}
	return fmt.Sprintf("<expr>(%s)", strings.Join(elt, " "))
}

func NewExpr(atom Node, list Node) *Expr {
	e := &Expr{}
	if atom != nil {
		e.children = append(e.children, atom)
	}
	if list != nil {
		e.children = append(e.children, list.Children()...)
	}

	return e
}

type Identifier struct {
	value string
}

func (i *Identifier) Children() []Node {
	return nil
}

func (i *Identifier) Eval(ctx Context) interface{} {
	return ctx.Get(i.value)
}

func (i *Identifier) String() string {
	return fmt.Sprintf("<id>%q", i.value)
}

func NewIdentifier(token Token) *Identifier {
	return &Identifier{
		value: token.Value().(string),
	}
}

type Integer struct {
	value int64
}

func (i *Integer) Children() []Node {
	return nil
}

func (i *Integer) Eval(ctx Context) interface{} {
	return i.value
}

func (i *Integer) String() string {
	return fmt.Sprintf("<integer>%d", i.value)
}

func NewInteger(token Token) *Integer {
	return &Integer{
		value: token.Value().(int64),
	}
}
