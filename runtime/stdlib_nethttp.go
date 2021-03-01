package runtime

import (
	"net/http"

	"github.com/rumlang/rum/parser"
)

// HTTPLib struct
type HTTPLib struct{}

// LoadLib function to HTTPLib struct
// https://gobyexample.com/http-servers
func (l *HTTPLib) LoadLib(ctx *Context, funcPrefix parser.Identifier) {
	if funcPrefix == "" {
		funcPrefix = "http"
	}
	ctx.SetFn(ConcatIdentifier(funcPrefix, ".listen-serve"), http.ListenAndServe, CheckArity(2))
	ctx.SetFn(ConcatIdentifier(funcPrefix, ".handle-func"), http.HandleFunc, CheckArity(2))
	// ctx.SetFn(ConcatIdentifier(funcPrefix, ".response-writer"), http.ResponseWriter, CheckArity(0))
	// ctx.SetFn(ConcatIdentifier(funcPrefix, ".request"), *http.Request, CheckArity(0))
}
