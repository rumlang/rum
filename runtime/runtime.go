package runtime

import (
	"fmt"
	"reflect"

	log "github.com/golang/glog"
	"github.com/palats/glop/parser"
)

const (
	ErrUnknownVariable = iota
)

type ErrorCode int

func (c ErrorCode) String() string {
	switch c {
	case ErrUnknownVariable:
		return "UnknownVariable"
	default:
		return fmt.Sprintf("Unknown[%d]", c)
	}
}

// Error is sent through panic when something went wrong during the execution.
type Error struct {
	Code ErrorCode
	Msg  string
	Ref  parser.SourceRef
}

func (e Error) String() string {
	return fmt.Sprintf("runtime error: %s[%d] - %s", e.Code, e.Code, e.Msg)
}

func (e Error) Error() string {
	return e.String()
}

// Internal is the type used to recognized internal functions (for which
// arguments are not evaluated automatically) from regular functions.
type Internal func(*Context, ...parser.Value) parser.Value

// Context contains details about the current execution frame.
type Context struct {
	parent *Context
	env    map[string]parser.Value
}

// Get returns the content of the specified variable. It will automatically
// look up parent context if needed. Generate a panic with an Error object if
// the specified variable does not exists.
func (c *Context) Get(s string) (parser.Value, error) {
	v, ok := c.env[s]
	if !ok {
		if c.parent != nil {
			return c.parent.Get(s)
		}
		return nil, Error{
			Code: ErrUnknownVariable,
			Msg:  fmt.Sprintf("%q does not exist", s),
		}
	}
	return v, nil
}

func (c *Context) Set(s string, v parser.Value) parser.Value {
	c.env[s] = v
	return v
}

// Eval takes the provided AST, evaluates it based on the current content of
// the context and returns the result.
// XXX should return a struct with metadata + value (a bit like Node in practice) for all type; i.e., always annotate what we're manipulating
func (c *Context) Eval(v parser.Value) parser.Value {
	node := v.(*parser.Node)
	switch node.Type {
	case parser.NodeExpression:
		log.Info("Expr:Eval", node)
		if len(node.Children()) <= 0 {
			return nil
		}

		fn := c.Eval(node.Children()[0].(*parser.Node)).Value()

		if internal, ok := fn.(Internal); ok {
			return internal(c, node.Children()[1:]...)
		}

		var args []reflect.Value
		for _, child := range node.Children()[1:] {
			args = append(args, reflect.ValueOf(c.Eval(child).Value()))
		}
		result := reflect.ValueOf(fn).Call(args)
		if len(result) == 0 {
			return parser.ValueAny(nil)
		}
		if len(result) == 1 {
			return parser.ValueAny(result[0].Interface())
		}
		panic("Multiple arguments unsupported")
	case parser.NodeIdentifier:
		v, err := c.Get(node.Value().(string))
		if err == nil {
			return v
		}
		if e, ok := err.(Error); ok {
			e.Ref = node.Ref
			panic(e)
		}
		panic(err)
	case parser.NodeInteger:
		return node
	case parser.NodeFloat:
		return node
	default:
		panic("EvalTODO")
	}
}

func NewContext(parent *Context) *Context {
	c := &Context{
		parent: parent,
		env:    make(map[string]parser.Value),
	}

	defaults := map[string]interface{}{
		"begin":  Begin,
		"quote":  Internal(Quote),
		"define": Internal(Define),
		"set!":   Internal(Define),
		"if":     Internal(If),
		"lambda": Internal(Lambda),

		"type": Type,

		"true":  true,
		"false": false,

		"+":        OpAdd,
		"+int64":   OpAddInt64,
		"+float64": OpAddFloat64,
		"-":        OpSub,
		"*":        OpMul,
		"*int64":   OpMulInt64,
		"*float64": OpMulFloat64,
		"==":       OpEqual,
		"eq?":      OpEqual,
		"!=":       OpNotEqual,
		"<":        OpLess,
		"<=":       OpLessEqual,
		">":        OpGreater,
		">=":       OpGreaterEqual,
	}

	for name, value := range defaults {
		c.env[name] = parser.ValueAny(value)
	}

	return c
}

func Quote(ctx *Context, args ...parser.Value) parser.Value {
	if len(args) != 1 {
		panic("Invalid number of arguments for quote")
	}
	return args[0]
}

func Begin(values ...interface{}) interface{} {
	if len(values) == 0 {
		return nil
	}
	return values[len(values)-1]
}

func Define(ctx *Context, args ...parser.Value) parser.Value {
	if len(args) != 2 {
		panic("Invalid arguments")
	}
	if args[0].(*parser.Node).Type != parser.NodeIdentifier {
		panic("TODO")
	}
	s := args[0].Value().(string)
	return ctx.Set(s, ctx.Eval(args[1]))
}

// If implements the 'if' builtin function. It has to be an Internal interface
// - otherwise, both true & false expressions would have been already
// evaluated.
func If(ctx *Context, args ...parser.Value) parser.Value {
	if len(args) < 2 || len(args) > 3 {
		panic("Invalid arguments")
	}

	cond := ctx.Eval(args[0]).Value().(bool)
	if cond {
		return ctx.Eval(args[1])
	}

	if len(args) <= 2 {
		return parser.ValueAny(nil)
	}

	return ctx.Eval(args[2])
}

func Lambda(ctx *Context, args ...parser.Value) parser.Value {
	if len(args) != 2 {
		panic("Invalid arguments")
	}

	names := []string{}
	for _, n := range args[0].(*parser.Node).Children() {
		if n.(*parser.Node).Type != parser.NodeIdentifier {
			panic("TODO")
		}
		names = append(names, n.Value().(string))
	}
	implNode := args[1]
	impl := func(implCtx *Context, args ...parser.Value) parser.Value {
		if len(args) != len(names) {
			panic("TODO")
		}
		nested := NewContext(implCtx)
		for i, name := range names {
			nested.Set(name, implCtx.Eval(args[i]))
		}
		return nested.Eval(implNode)
	}

	return parser.ValueAny(Internal(impl))
}

func Type(v interface{}) string {
	return fmt.Sprintf("%T", v)
}

func ParseEval(src *parser.Source) (interface{}, []error) {
	node, errs := parser.Parse(src)
	if len(errs) > 0 {
		return nil, errs
	}
	return NewContext(nil).Eval(node).Value(), nil
}
