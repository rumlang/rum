package runtime

import (
	"github.com/palats/glop/nodes"
	"github.com/palats/glop/parser"
)

type Context struct {
	parent nodes.Context
	env    map[string]interface{}
}

func (c *Context) Get(s string) interface{} {
	v, ok := c.env[s]
	if !ok && c.parent != nil {
		return c.parent.Get(s)
	}
	return v
}

func (c *Context) Set(s string, v interface{}) interface{} {
	c.env[s] = v
	return v
}

func NewContext(parent nodes.Context) *Context {
	c := &Context{
		parent: parent,
		env:    make(map[string]interface{}),
	}

	c.env["begin"] = Begin
	c.env["quote"] = nodes.Internal(Quote)
	c.env["define"] = nodes.Internal(Define)
	c.env["set!"] = nodes.Internal(Define)
	c.env["if"] = nodes.Internal(If)
	c.env["lambda"] = nodes.Internal(Lambda)

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

func Quote(ctx nodes.Context, args ...nodes.Node) interface{} {
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

func Define(ctx nodes.Context, args ...nodes.Node) interface{} {
	if len(args) != 2 {
		panic("Invalid arguments")
	}
	s := args[0].(nodes.Identifier).Value()
	return ctx.Set(s, args[1].Eval(ctx))
}

// If implements the 'if' builtin function. It has to be an Internal interface
// - otherwise, both true & false expressions would have been already
// evaluated.
func If(ctx nodes.Context, args ...nodes.Node) interface{} {
	if len(args) < 2 || len(args) > 3 {
		panic("Invalid arguments")
	}

	cond := args[0].Eval(ctx).(bool)
	if cond {
		return args[1].Eval(ctx)
	}

	if len(args) <= 2 {
		return nil
	}

	return args[2].Eval(ctx)
}

func Lambda(ctx nodes.Context, args ...nodes.Node) interface{} {
	if len(args) != 2 {
		panic("Invalid arguments")
	}

	names := []string{}
	for _, n := range args[0].Children() {
		names = append(names, n.(nodes.Identifier).Value())
	}
	implNode := args[1]
	impl := func(ctx nodes.Context, args ...nodes.Node) interface{} {
		if len(args) != len(names) {
			panic("TODO")
		}
		nested := NewContext(ctx)
		for i, name := range names {
			nested.Set(name, args[i].Eval(ctx))
		}
		return implNode.Eval(nested)
	}

	return nodes.Internal(impl)
}

func ParseEval(src *parser.Source) (interface{}, []error) {
	node, errs := parser.Parse(src)
	if len(errs) > 0 {
		return nil, errs
	}
	return node.Eval(NewContext(nil)), nil
}
