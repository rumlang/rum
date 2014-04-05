package runner

import (
	"reflect"
	"testing"

	"github.com/palats/glop/nodes"
)

func TestEval(t *testing.T) {
	r := ParseEval("(+ 1 2)")
	if r.(int64) != 3 {
		t.Errorf("Expected a result of 3, got: %v", r)
	}
}

func TestQuote(t *testing.T) {
	n := ParseEval("(quote (+ 1 2))").(nodes.Node)

	if len(n.Children()) != 3 {
		t.Errorf("Expected 3 children, got: %v", n.Children())
	}
}

func TestBegin(t *testing.T) {
	n := ParseEval("(begin 1 (+ 1 1))")
	if n.(int64) != 2 {
		t.Errorf("Expected '2', got: %v", n)
	}

	n = ParseEval("(begin)")
	if n != nil {
		t.Errorf("Expected nil, got: %v", n)
	}
}

func TestDefine(t *testing.T) {
	n := ParseEval("(begin (define a 5) a)")
	if n.(int64) != 5 {
		t.Errorf("Expected '5', got: %v", n)
	}
}

func TestSet(t *testing.T) {
	n := ParseEval("(begin (define a 5) (set! a 4) a)")
	if n.(int64) != 4 {
		t.Errorf("Expected '4', got: %v", n)
	}
}

func TestIf(t *testing.T) {
	valid := map[string]interface{}{
		"(if true 7)":    int64(7),
		"(if false 7)":   nil,
		"(if false 7 8)": int64(8),
	}

	for input, expected := range valid {
		r := ParseEval(input)
		if !reflect.DeepEqual(r, expected) {
			t.Errorf("Input %q -- expected <%T>%#+v, got: <%T>%#+v", input, expected, expected, r, r)
		}
	}
}
