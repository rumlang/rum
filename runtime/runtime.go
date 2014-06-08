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
	Ref  *parser.SourceRef
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
	env    map[parser.Identifier]parser.Value
}

// Get returns the content of the specified variable. It will automatically
// look up parent context if needed. Generate a panic with an Error object if
// the specified variable does not exists.
func (c *Context) Get(id parser.Identifier) (parser.Value, error) {
	v, ok := c.env[id]
	if !ok {
		if c.parent != nil {
			return c.parent.Get(id)
		}
		return nil, Error{
			Code: ErrUnknownVariable,
			Msg:  fmt.Sprintf("%q does not exist", string(id)),
		}
	}
	return v, nil
}

func (c *Context) Set(id parser.Identifier, v parser.Value) parser.Value {
	c.env[id] = v
	return v
}

// Eval takes the provided value, evaluates it based on the current content of
// the context and returns the result.
func (c *Context) Eval(input parser.Value) parser.Value {
	switch data := input.Value().(type) {
	case []parser.Value:
		log.Info("Expr:Eval", input)

		if len(data) <= 0 {
			return nil
		}

		fn := c.Eval(data[0]).Value()

		if internal, ok := fn.(Internal); ok {
			return internal(c, data[1:]...)
		}

		var args []reflect.Value
		for _, child := range data[1:] {
			args = append(args, reflect.ValueOf(c.Eval(child).Value()))
		}
		result := reflect.ValueOf(fn).Call(args)
		if len(result) == 0 {
			return parser.NewAny(nil, nil)
		}
		if len(result) == 1 {
			return parser.NewAny(result[0].Interface(), nil)
		}
		panic("Multiple arguments unsupported")
	case parser.Identifier:
		v, err := c.Get(data)
		if err == nil {
			return v
		}
		if e, ok := err.(Error); ok {
			e.Ref = input.Ref()
			panic(e)
		}
		panic(err)
	default:
		// If it is neither an identifier or a list, just return the value.
		return input
	}
}

func NewContext(parent *Context) *Context {
	c := &Context{
		parent: parent,
		env:    make(map[parser.Identifier]parser.Value),
	}

	defaults := map[parser.Identifier]interface{}{
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
		c.env[name] = parser.NewAny(value, nil)
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

	id, ok := args[0].Value().(parser.Identifier)
	if !ok {
		panic("TODO")
	}
	return ctx.Set(id, ctx.Eval(args[1]))
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
		return parser.NewAny(nil, nil)
	}

	return ctx.Eval(args[2])
}

func Lambda(ctx *Context, args ...parser.Value) parser.Value {
	if len(args) != 2 {
		panic("Invalid arguments")
	}

	params, ok := args[0].Value().([]parser.Value)
	if !ok {
		panic("TODO")
	}
	names := []parser.Identifier{}
	for _, v := range params {
		id, ok := v.Value().(parser.Identifier)
		if !ok {
			panic("TODO")
		}
		names = append(names, id)
	}
	implValue := args[1]
	impl := func(implCtx *Context, args ...parser.Value) parser.Value {
		if len(args) != len(names) {
			panic("TODO")
		}
		nested := NewContext(implCtx)
		for i, name := range names {
			nested.Set(name, implCtx.Eval(args[i]))
		}
		return nested.Eval(implValue)
	}

	return parser.NewAny(Internal(impl), nil)
}

func Type(v interface{}) string {
	return fmt.Sprintf("%T", v)
}

func ParseEval(src *parser.Source) (interface{}, []error) {
	v, errs := parser.Parse(src)
	if len(errs) > 0 {
		return nil, errs
	}
	return NewContext(nil).Eval(v).Value(), nil
}
