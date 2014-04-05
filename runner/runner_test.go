package runner

import (
	"reflect"
	"testing"

	"github.com/palats/glop/nodes"
)

func TestQuote(t *testing.T) {
	n := ParseEval("(quote (+ 1 2))").(nodes.Node)

	if len(n.Children()) != 3 {
		t.Errorf("Expected 3 children, got: %v", n.Children())
	}
}

func TestValid(t *testing.T) {
	valid := map[string]interface{}{
		"(+ 1 2)":                           int64(3),
		"(begin 1 (+ 1 1))":                 int64(2),
		"(begin)":                           nil,
		"(begin (define a 5) a)":            int64(5),
		"(begin (define a 5) (set! a 4) a)": int64(4),
		"(if true 7)":                       int64(7),
		"(if false 7)":                      nil,
		"(if false 7 8)":                    int64(8),
	}

	for input, expected := range valid {
		r := ParseEval(input)
		if !reflect.DeepEqual(r, expected) {
			t.Errorf("Input %q -- expected <%T>%#+v, got: <%T>%#+v", input, expected, expected, r, r)
		}
	}
}
