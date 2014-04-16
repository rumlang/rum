// Package node provide the core implementation of the nodes of the AST. It
// provides the implementation of evaluation for each of those.
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

type Expr []Node

func (e Expr) Children() []Node {
	return e
}

func (e Expr) Eval(ctx Context) interface{} {
	log.Info("Expr:Eval", e)
	if len(e.Children()) <= 0 {
		return nil
	}

	fn := e.Children()[0].Eval(ctx)

	if internal, ok := fn.(Internal); ok {
		return internal(ctx, e.Children()[1:]...)
	}

	var args []reflect.Value
	for _, children := range e.Children()[1:] {
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

func (e Expr) String() string {
	var elt []string
	for _, node := range e.Children() {
		elt = append(elt, node.String())
	}
	return fmt.Sprintf("<expr>(%s)", strings.Join(elt, " "))
}

func NewExpr(atoms []Node) Expr {
	return Expr(atoms)
}

type Identifier string

func (i Identifier) Children() []Node {
	return nil
}

func (i Identifier) Eval(ctx Context) interface{} {
	return ctx.Get(string(i))
}

func (i Identifier) String() string {
	return fmt.Sprintf("<id>%q", string(i))
}

func (i Identifier) Value() string {
	return string(i)
}

func NewIdentifier(token Token) Identifier {
	return Identifier(token.Value().(string))
}

type Integer int64

func (i Integer) Children() []Node {
	return nil
}

func (i Integer) Eval(ctx Context) interface{} {
	return int64(i)
}

func (i Integer) String() string {
	return fmt.Sprintf("<integer>%d", int64(i))
}

func NewInteger(token Token) Integer {
	return Integer(token.Value().(int64))
}
