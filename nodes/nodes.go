package nodes


type Node interface {
  Raw() string
  Children() []Node
}

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

func NewIdentifier(raw string) *Identifier {
  return &Identifier{
    raw: raw,
  }
}
