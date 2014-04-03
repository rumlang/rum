package runner

import (
	"testing"

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
