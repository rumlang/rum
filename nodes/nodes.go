package nodes

import (
  "fmt"
  "log"
  "strings"
)

type Node interface {
  Raw() string
  Children() []Node
  String() string
  Exec(ctx Context) interface{}
}

type Context interface {
  Get(s string) interface{}
}

type Internal func(Context, ...Node) interface{}

type Expr struct {
  raw string
  children []Node
}

func (e *Expr) Raw() string {
  return e.raw
}

func (e *Expr) Children() []Node {
  return e.children
}

func (e *Expr) Exec(ctx Context) interface{} {
  log.Print("Expr:Exec", e)
  if len(e.children) <= 0 {
    return nil
  }

  fn := e.children[0].Exec(ctx)
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
  raw string
}

func (i *Identifier) Raw() string {
  return i.raw
}

func (i *Identifier) Children() []Node {
  return nil
}

func (i *Identifier) Exec(ctx Context) interface{} {
  return ctx.Get(i.raw)
}

func (i *Identifier) String() string {
  return fmt.Sprintf("<id>%q", i.raw)
}

func NewIdentifier(raw string) *Identifier {
  return &Identifier{
    raw: raw,
  }
}
