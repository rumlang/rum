package parser

import (
	"testing"
)

func TestLexer(t *testing.T) {
	tests := map[string][]tokenInfo{
		"foo": []tokenInfo{
			{text: "foo", id: tokIdentifier, value: "foo"},
		},
		"(foo)": []tokenInfo{
			{text: "(", id: tokOpen},
			{text: "foo", id: tokIdentifier, value: "foo"},
			{text: ")", id: tokClose},
		},
		" (  foo ) ": []tokenInfo{
			{text: "(", id: tokOpen},
			{text: "foo", id: tokIdentifier, value: "foo"},
			{text: ")", id: tokClose},
		},
	}

	for input, expected := range tests {
		l := newLexer(input)

		tokens := []tokenInfo{}
		for t := range l.tokens {
			tokens = append(tokens, t)
		}
		if len(tokens) != len(expected) {
			t.Fatalf("Invalid number of tokens found; expected %v, found %v", expected, tokens)
		}
		for i, token := range tokens {
			if token != expected[i] {
				t.Errorf("Expression %q - expected %v, got %v", input, expected[i], token)
			}
		}
	}
}

func TestParsing(t *testing.T) {
	tests := map[string]int{
		"foo":               0,
		"a(b":               -1,
		"()":                0,
		"(foo)":             1,
		"(a b)":             2,
		"(a (b c))":         2,
		"(a (b c) d (e f))": 4,
	}

	for input, count := range tests {
		l := newLexer(input)
		yyParse(l)

		if count < 0 {
			if len(l.errors) == 0 {
				t.Errorf("Input %q parsed instead of generating error", input)
			}
			continue
		}

		if len(l.errors) != 0 {
			t.Errorf("Input %q - parsing errors: %v", input, l.errors)
		}

		if len(l.program.Children()) != count {
			t.Errorf("Input %q - expecting %d children, got %d: %v", input, count, len(l.program.Children()), l.program.Children())
		}
	}
}
