package parser

import (
	"reflect"
	"testing"
)

func TestLexer(t *testing.T) {
	tests := map[string][]tokenInfo{
		"foo": []tokenInfo{
			{text: []rune{'f', 'o', 'o'}, id: tokIdentifier, value: "foo", ref: SourceRef{Line: 0, Column: 0}},
		},
		"(foo)": []tokenInfo{
			{text: []rune{'('}, id: tokOpen, ref: SourceRef{Line: 0, Column: 0}},
			{text: []rune{'f', 'o', 'o'}, id: tokIdentifier, value: "foo", ref: SourceRef{Line: 0, Column: 1}},
			{text: []rune{')'}, id: tokClose, ref: SourceRef{Line: 0, Column: 4}},
		},
		" (  foo ) ": []tokenInfo{
			{text: []rune{'('}, id: tokOpen, ref: SourceRef{Line: 0, Column: 1}},
			{text: []rune{'f', 'o', 'o'}, id: tokIdentifier, value: "foo", ref: SourceRef{Line: 0, Column: 4}},
			{text: []rune{')'}, id: tokClose, ref: SourceRef{Line: 0, Column: 8}},
		},
		" \nfoo": []tokenInfo{
			{text: []rune{'f', 'o', 'o'}, id: tokIdentifier, value: "foo", ref: SourceRef{Line: 1, Column: 0}},
		},
	}

	for input, expected := range tests {
		l := newLexer(input)

		tokens := []tokenInfo{}
		for t := range l.tokens {
			tokens = append(tokens, t)
		}
		if len(tokens) != len(expected) {
			t.Fatalf("Expression %q - invalid number of tokens found; expected %#+v, found %#+v", input, expected, tokens)
		}
		for i, token := range tokens {
			if !reflect.DeepEqual(token, expected[i]) {
				t.Errorf("Expression %q - expected %#+v, got %#+v", input, expected[i], token)
			}
		}
	}
}

func TestLexerErrors(t *testing.T) {
	tests := map[string]int{
		// Invalid sequence at the beginning - skips the first byte and so gets a
		// ')'.
		"\xc3\x28":   1,
		"a\xc3\x28b": 3,
	}

	for input, count := range tests {
		l := newLexer(input)
		for _ = range l.tokens {
			count--
		}
		if count > 0 {
			t.Errorf("Input %# x ; Not enough token found - %d more expected", input, count)
		}
		if count < 0 {
			t.Errorf("Input %# x ; Too many token found - %d more than expected", input, -count)
		}
	}
}

func TestParsing(t *testing.T) {
	tests := map[string]int{
		"":                  -1,
		")":                 -1,
		"(":                 -1,
		"a)b":               -1,
		")b":                -1,
		"a(b":               -1,
		"()":                0,
		"foo":               0,
		"(foo)":             1,
		"(a b)":             2,
		"(a (b c))":         2,
		"(a (b c) d (e f))": 4,
		"(a\nb)":            2,
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

func TestParsingErrors(t *testing.T) {
	type foo struct {
		code ErrorCode
		ref  SourceRef
	}

	tests := map[string]foo{
		"(+": foo{
			code: ErrMissingClosingParenthesis,
			ref: SourceRef{
				Line:   0,
				Column: 2,
			},
		},
		"(": foo{
			code: ErrMissingClosingParenthesis,
			ref: SourceRef{
				Line:   0,
				Column: 1,
			},
		},
	}

	for input, expected := range tests {
		_, errs := Parse(input)
		if len(errs) != 1 {
			t.Errorf("Input %q should have 1 error; instead: %v", input, errs)
		} else {
			err := errs[0].(Error)
			if err.Code != expected.code {
				t.Errorf("Input %q should have returned error code %d; instead: %d", input, expected.code, err.Code)
			}
			if !reflect.DeepEqual(err.Ref, expected.ref) {
				t.Errorf("Input %q - expected %#+v, got %#+v", input, expected.ref, err.Ref)
			}
		}
	}
}
