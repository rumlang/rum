package runtime

import (
	"reflect"
	"testing"

	"github.com/palats/glop/parser"
)

func parseEval(src *parser.Source) (interface{}, error) {
	v, err := parser.Parse(src)
	if err != nil {
		return nil, err
	}
	return NewContext(nil).Eval(v).Value(), nil
}

func TestQuote(t *testing.T) {
	r, err := parseEval(parser.NewSource("(quote (+ 1 2))"))
	if err != nil {
		t.Fatal(err)
	}
	n := r.([]parser.Value)

	if len(n) != 3 {
		t.Errorf("Expected 3 children, got: %v", n)
	}
}

func TestValid(t *testing.T) {
	valid := map[string]interface{}{
		// Test '+'
		"(+ 1 2)":            int64(3),
		"(+int64 1 2)":       int64(3),
		"(+ 1.0 2.0)":        float64(3),
		"(+float64 1.0 2.0)": float64(3),
		"(- 1 2)":            int64(-1),
		"(- 1.0 2.0)":        float64(-1.0),
		"(* 3.0 2.0)":        float64(6),
		"(* 3 2)":            int64(6),
		"(== 3 2)":           false,
		"(== 3 3)":           true,
		"(== 3 3 1)":         false,
		"(== 3 3 3)":         true,
		"(== 3.0 3.0)":       true,
		"(== 3.0 2.0)":       false,
		"(!= 3 2)":           true,
		"(!= 3 3)":           false,
		"(!= 3.0 3.0)":       false,
		"(!= 3.0 2.0)":       true,
		"(< 3 2)":            false,
		"(< 2 3)":            true,
		"(< 3.0 3.0)":        false,
		"(< 3.0 2.0)":        false,
		"(<= 3 2)":           false,
		"(<= 2 3)":           true,
		"(<= 3.0 3.0)":       true,
		"(<= 3.0 2.0)":       false,
		"(> 3 2)":            true,
		"(> 2 3)":            false,
		"(> 3.0 3.0)":        false,
		"(> 3.0 2.0)":        true,
		"(>= 3 2)":           true,
		"(>= 2 3)":           false,
		"(>= 3.0 3.0)":       true,
		"(>= 3.0 2.0)":       true,
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
		// Test float
		".3": float64(.3),
		// Test car
		"(car (quote (1 2 3)))": int64(1),
		// Test length
		"(length (quote (1 2 3)))": int64(3),
	}

	for input, expected := range valid {
		r, err := parseEval(parser.NewSource(input))
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(r, expected) {
			t.Errorf("Input %q -- expected <%T>%#+v, got: <%T>%#+v", input, expected, expected, r, r)
		}
	}
}

func TestValidList(t *testing.T) {
	valid := map[string][]interface{}{
		// Test cons
		"(cons 1 (quote (2 3)))":                                    []interface{}{int64(1), int64(2), int64(3)},
		"(begin (define a (quote (1 2 3))) (cons (car a) (cdr a)))": []interface{}{int64(1), int64(2), int64(3)},
		// Test cdr
		"(cdr (quote (1 2 3)))": []interface{}{int64(2), int64(3)},
	}

	for input, expected := range valid {
		r, err := parseEval(parser.NewSource(input))
		if err != nil {
			t.Fatal(err)
		}

		if _, ok := r.([]parser.Value); !ok {
			t.Errorf("Expected a []Value; got: %T", r)
		}

		for i, v := range r.([]parser.Value) {
			if v.Value() != expected[i] {
				t.Errorf("Item %d - expected %v, got: %v", i, expected[i], v.Value())
			}
		}
	}
}

func TestPanic(t *testing.T) {
	panics := []string{
		"(6)",
		"(+ 1 (2))",
		"(+ 1.0 2)",
		"(+int64 1.0 2.0)",
		"(+float64 1 2)",
		"(*int64 1.0 2.0)",
		"(*float64 1 2)",
		"(cons 1 2)",
		"(cons 1 (2 3))",
	}

	for _, s := range panics {
		var r interface{}
		func() {
			defer func() {
				r = recover()
			}()
			parseEval(parser.NewSource(s))
		}()

		if r == nil {
			t.Fatalf("%q should have generated a panic.", s)
		}
	}
}

func TestUnknownVariable(t *testing.T) {
	var r interface{}
	func() {
		defer func() {
			r = recover()
		}()
		parseEval(parser.NewSource("(a)"))
	}()
	if r == nil {
		t.Fatalf("Expected a panic, got nothing")
	}

	e, ok := r.(Error)
	if !ok {
		t.Fatalf("Expected a runtime.Error; instead: %v", r)
	}
	if e.Code != ErrUnknownVariable {
		t.Errorf("Expected an UnknownVariable; instead: %v", e)
	}
}
