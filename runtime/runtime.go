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
}

func (e Error) String() string {
	return fmt.Sprintf("runtime error: %s[%d] - %s", e.Code, e.Code, e.Msg)
}

// Internal is the type used to recognized internal functions (for which
// arguments are not evaluated automatically) from regular functions.
type Internal func(*Context, ...*parser.Node) interface{}

// Context contains details about the current execution frame.
type Context struct {
	parent *Context
	env    map[string]interface{}
}

// Get returns the content of the specified variable. It will automatically
// look up parent context if needed. Generate a panic with an Error object if
// the specified variable does not exists.
func (c *Context) Get(s string) interface{} {
	v, ok := c.env[s]
	if !ok {
		if c.parent != nil {
			return c.parent.Get(s)
		}
		panic(Error{
			Code: ErrUnknownVariable,
			Msg:  fmt.Sprintf("%q does not exists", s),
		})
	}
	return v
}

func (c *Context) Set(s string, v interface{}) interface{} {
	c.env[s] = v
	return v
}

// Eval takes the provided AST, evaluates it based on the current content of
// the context and returns the result.
func (c *Context) Eval(node *parser.Node) interface{} {
	switch node.Type {
	case parser.NodeExpression:
		log.Info("Expr:Eval", node)
		if len(node.Children()) <= 0 {
			return nil
		}

		fn := c.Eval(node.Children()[0])

		if internal, ok := fn.(Internal); ok {
			return internal(c, node.Children()[1:]...)
		}

		var args []reflect.Value
		for _, child := range node.Children()[1:] {
			args = append(args, reflect.ValueOf(c.Eval(child)))
		}
		result := reflect.ValueOf(fn).Call(args)
		if len(result) == 0 {
			return nil
		}
		if len(result) == 1 {
			return result[0].Interface()
		}
		panic("Multiple arguments unsupported")
	case parser.NodeIdentifier:
		return c.Get(node.Value.(string))
	case parser.NodeInteger:
		return node.Value.(int64)
	case parser.NodeFloat:
		return node.Value.(float64)
	default:
		panic("EvalTODO")
	}
}

func NewContext(parent *Context) *Context {
	c := &Context{
		parent: parent,
		env:    make(map[string]interface{}),
	}

	c.env["begin"] = Begin
	c.env["quote"] = Internal(Quote)
	c.env["define"] = Internal(Define)
	c.env["set!"] = Internal(Define)
	c.env["if"] = Internal(If)
	c.env["lambda"] = Internal(Lambda)

	c.env["+"] = OpAdd
	c.env["true"] = true
	c.env["false"] = false

	return c
}

func OpAdd(values ...int64) int64 {
	var total int64
	for _, v := range values {
		total += v
	}
	return total
}

func Quote(ctx *Context, args ...*parser.Node) interface{} {
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

func Define(ctx *Context, args ...*parser.Node) interface{} {
	if len(args) != 2 {
		panic("Invalid arguments")
	}
	if args[0].Type != parser.NodeIdentifier {
		panic("TODO")
	}
	s := args[0].Value.(string)
	return ctx.Set(s, ctx.Eval(args[1]))
}

// If implements the 'if' builtin function. It has to be an Internal interface
// - otherwise, both true & false expressions would have been already
// evaluated.
func If(ctx *Context, args ...*parser.Node) interface{} {
	if len(args) < 2 || len(args) > 3 {
		panic("Invalid arguments")
	}

	cond := ctx.Eval(args[0]).(bool)
	if cond {
		return ctx.Eval(args[1])
	}

	if len(args) <= 2 {
		return nil
	}

	return ctx.Eval(args[2])
}

func Lambda(ctx *Context, args ...*parser.Node) interface{} {
	if len(args) != 2 {
		panic("Invalid arguments")
	}

	names := []string{}
	for _, n := range args[0].Children() {
		if n.Type != parser.NodeIdentifier {
			panic("TODO")
		}
		names = append(names, n.Value.(string))
	}
	implNode := args[1]
	impl := func(implCtx *Context, args ...*parser.Node) interface{} {
		if len(args) != len(names) {
			panic("TODO")
		}
		nested := NewContext(implCtx)
		for i, name := range names {
			nested.Set(name, implCtx.Eval(args[i]))
		}
		return nested.Eval(implNode)
	}

	return Internal(impl)
}

func ParseEval(src *parser.Source) (interface{}, []error) {
	node, errs := parser.Parse(src)
	if len(errs) > 0 {
		return nil, errs
	}
	return NewContext(nil).Eval(node), nil
}
