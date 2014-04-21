package runtime

import (
	"reflect"
	"testing"

	"github.com/palats/glop/nodes"
)

func TestQuote(t *testing.T) {
	r, err := ParseEval("(quote (+ 1 2))")
	if err != nil {
		t.Fatal(err)
	}
	n := r.(nodes.Node)

	if len(n.Children()) != 3 {
		t.Errorf("Expected 3 children, got: %v", n.Children())
	}
}

func TestValid(t *testing.T) {
	valid := map[string]interface{}{
		// Test '+'
		"(+ 1 2)": int64(3),
		// Test 'begin'
		"(begin 1 (+ 1 1))": int64(2),
		"(begin)":           nil,
		// Test 'define'
		"(begin (define a 5) a)":            int64(5),
		"(begin (define a 5) (set! a 4) a)": int64(4),
		// Test 'if'
		"(if true 7)":    int64(7),
		"(if false 7)":   nil,
		"(if false 7 8)": int64(8),
		// Test 'lambda'
		"(begin (define d (lambda (n) (+ n n))) (d 3))": int64(6),
		// Test that inner scopes are not override outer scope.
		"(begin (define n 7) (define d (lambda (n) (+ n n))) (+ n (d 3)))": int64(13),
	}

	for input, expected := range valid {
		r, err := ParseEval(input)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(r, expected) {
			t.Errorf("Input %q -- expected <%T>%#+v, got: <%T>%#+v", input, expected, expected, r, r)
		}
	}
}

func TestPanic(t *testing.T) {
	panics := []string{
		"(6)",
		"(+ 1 (2))",
	}

	for _, s := range panics {
		var r interface{}
		func() {
			defer func() {
				r = recover()
			}()
			ParseEval("(6)")
		}()

		if r == nil {
			t.Fatalf("%q should have generated a panic.", s)
		}
	}
}
