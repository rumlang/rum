package parser

import "testing"

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
  tests := []string{
    "foo",
    "()",
    "(foo)",
    "(a (b c) d (e f))",
  }

  for _, input := range tests {
    l := newLexer(input)
    yyParse(l)
    if len(l.errors) != 0 {
      t.Fatalf("Input %q - parsing errors: %v", input, l.errors)
    }
  }
}
