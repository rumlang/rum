package nodes

import (
	"fmt"
	"reflect"
	"strings"

	log "github.com/golang/glog"
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
	Set(s string, v interface{}) interface{}
}

type Internal func(Context, ...Node) interface{}

type Expr struct {
	children []Node
}

func (e *Expr) Children() []Node {
	return e.children
}

func (e *Expr) Eval(ctx Context) interface{} {
	log.Info("Expr:Eval", e)
	if len(e.children) <= 0 {
		return nil
	}

	fn := e.children[0].Eval(ctx)

	if internal, ok := fn.(Internal); ok {
		return internal(ctx, e.children[1:]...)
	}

	var args []reflect.Value
	for _, children := range e.children[1:] {
		args = append(args, reflect.ValueOf(children.Eval(ctx)))
	}
	result := reflect.ValueOf(fn).Call(args)
	if len(result) == 0 {
		return nil
	}
	if len(result) == 1 {
		return result[0].Interface()
	}
	panic("Multiple arguments unsupportted")
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

func (i *Identifier) Value() string {
	return i.value
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
