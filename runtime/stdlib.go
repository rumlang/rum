package runtime

import (
	"encoding/csv"
	"fmt"
	"strings"
)

type StdLib interface {
	LoadLib(ctx *Context)
}

// StringLib struct
type StringLib struct{}

// LoadLib function to StringLib struct
func (l *StringLib) LoadLib(ctx *Context) {
	ctx.SetFn("strings.Compare", strings.Compare, CheckArity(2))
	ctx.SetFn("strings.Contains", strings.Contains, CheckArity(2))
	ctx.SetFn("strings.Count", strings.Count, CheckArity(2))
	ctx.SetFn("strings.Join", strings.Join, CheckArity(2))
	ctx.SetFn("strings.Split", strings.Split, CheckArity(2))
	ctx.SetFn("strings.Title", strings.Title, CheckArity(1))
	ctx.SetFn("strings.ToLower", strings.ToLower, CheckArity(1))
	ctx.SetFn("strings.ToUpper", strings.ToUpper, CheckArity(1))
	ctx.SetFn("strings.Trim", strings.Trim, CheckArity(2))
	ctx.SetFn("strings.NewReader", strings.NewReader, CheckArity(1))
}

// CSVLib struct
type CSVLib struct{}

// LoadLib function to StringLib struct
func (l *CSVLib) LoadLib(ctx *Context) {
	ctx.SetFn("encoding/csv.NewReader", csv.NewReader, CheckArity(1))
	ctx.SetFn("encoding/csv.NewWriter", csv.NewWriter, CheckArity(1))
}

func loadStdLib(name string, ctx *Context) {
	var stdLib StdLib
	switch name {
	case "strings":
		stdLib = &StringLib{}
		stdLib.LoadLib(ctx)
		return
	case "csv":
		stdLib = &CSVLib{}
		stdLib.LoadLib(ctx)
		return
	default:
		panic(fmt.Sprintf("package %s not found", name))
	}
}
