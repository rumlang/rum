package runner

import (
	"testing"

	"github.com/palats/glop/parser"
)

func TestExec(t *testing.T) {
	n := parser.Parse("(+ 1 2)")
	ctx := NewContext()
	r := n.Exec(ctx)
	if r.(int64) != 3 {
		t.Errorf("Expected a result of 3, got: %v", r)
	}
}
