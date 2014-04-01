package parser

import (
  "fmt"
  "reflect"
  "testing"
)

func TestLexer(t *testing.T) {
  tests := map[string][]tokenInfo{
    "foo": []tokenInfo{
      {"foo", tokIdentifier},
    },
    "(foo)": []tokenInfo{
      {"(", tokOpen},
      {"foo", tokIdentifier},
      {")", tokClose},
    },
    " (  foo ) ": []tokenInfo{
      {"(", tokOpen},
      {"foo", tokIdentifier},
      {")", tokClose},
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
  tests := map[string]bool{
    "foo": true,
    "a(b": false,
    "()": true,
    "(foo)": true,
    "(a b)": true,
    "(a (b c))": true,
    "(a (b c) d (e f))": true,
  }

  for input, valid := range tests {
    l := newLexer(input)
    yyParse(l)

    if valid && len(l.errors) != 0 {
      t.Errorf("Input %q - parsing errors: %v", input, l.errors)
    }

    if !valid && len(l.errors) == 0 {
      t.Errorf("Input %q parsed instead of generating error", input)
    }
  }
}
