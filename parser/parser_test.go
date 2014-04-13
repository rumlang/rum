package parser

import "testing"

func TestLexer(t *testing.T) {
	tests := map[string][]tokenInfo{
		"foo": []tokenInfo{
			{text: "foo", start: 0, id: tokIdentifier, value: "foo"},
		},
		"(foo)": []tokenInfo{
			{text: "(", start: 0, id: tokOpen},
			{text: "foo", start: 1, id: tokIdentifier, value: "foo"},
			{text: ")", start: 4, id: tokClose},
		},
		" (  foo ) ": []tokenInfo{
			{text: "(", start: 1, id: tokOpen},
			{text: "foo", start: 4, id: tokIdentifier, value: "foo"},
			{text: ")", start: 8, id: tokClose},
		},
		"\nfoo": []tokenInfo{
			{text: "foo", start: 1, line: 1, id: tokIdentifier, value: "foo"},
		},
	}

	for input, expected := range tests {
		l := newLexer(input)

		tokens := []tokenInfo{}
		for t := range l.tokens {
			tokens = append(tokens, t)
		}
		if len(tokens) != len(expected) {
			t.Fatalf("Invalid number of tokens found; expected %+v, found %+v", expected, tokens)
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
		"(a\nb)":            2,
		"a)b":               -1,
		")b":                -1,
		// ")":                 -1,
		// "(": -1,
	}

	for input, count := range tests {
		r, err := Parse(input)

		if count < 0 {
			if err == nil {
				t.Errorf("Input %q parsed instead of generating error", input)
			}
			continue
		}

		if err != nil {
			t.Errorf("Input %q - parsing errors: %v", input, err)
			continue
		}

		if count != len(r.Children()) {
			t.Errorf("Input %q - expected %d children, got %d: %v", input, count, len(r.Children()), r)
		}
	}
}
