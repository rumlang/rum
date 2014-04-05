package runner

import (
	"testing"

	"github.com/palats/glop/nodes"
	"github.com/palats/glop/parser"
)

func TestEval(t *testing.T) {
	n := parser.Parse("(+ 1 2)")
	ctx := NewContext()
	r := n.Eval(ctx)
	if r.(int64) != 3 {
		t.Errorf("Expected a result of 3, got: %v", r)
	}
}

func TestQuote(t *testing.T) {
	n := parser.Parse("(quote (+ 1 2))")
	ctx := NewContext()
	n = n.Eval(ctx).(nodes.Node)

	if len(n.Children()) != 3 {
		t.Errorf("Expected 3 children, got: %v", n.Children())
	}
}
