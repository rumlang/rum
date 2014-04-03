package runner

import (
	"github.com/palats/glop/nodes"
)

type Context struct {
	env map[string]interface{}
}

func (c *Context) Get(s string) interface{} {
	return c.env[s]
}

func NewContext() *Context {
	c := &Context{
		env: make(map[string]interface{}),
	}

	c.env["+"] = nodes.Internal(func(ctx nodes.Context, args ...nodes.Node) interface{} {
		var v int64
		for _, n := range args {
			v += n.Eval(ctx).(int64)
		}
		return v
	})

	return c
}
