package parser

import (
  "log"
  "strings"
)

type stateFn func() stateFn

type item struct {
  raw string
  id int
}

type lexer struct {
  input string
  start int
  pos int
  items chan item
}

func (l *lexer) peek() rune {
  r, _ := utf8.DecodeRuneInString(l.input[l.pos:])
  // XXX error detection
  return r
}

func (l *lexer) advance() {
  r, w := utf8.DecodeRuneInString(l.input[l.pos:])
  l.pos += w
}

func (l *lexer) emit(token int) {
  s := l.input[l.start:l.pos]
  // XXX emit
  l.start = l.pos
}

func (l *lexer) stateIdentifier() stateFn {
  identifier := []rune{}
  for {
    r := l.peek()

    switch {
    case r == '(':
      l.emit(tokIdentifier, identifier)
      l.advance()
      identifier = []rune{}
      l.emit(tokOpen, []rune{r})
    case r == ')':
      l.emit(tokIdentifier, identifier)
      l.advance()
      identifier = []rune{}
      l.emit(tokClose, []rune{r})
    case unicode.IsSpace(r):
      l.emit(tokIdentifier)
      return l.stateSpace
    default:
      identifier = append(identifier, r)
      l.advance()
    }
  }
}

func (l *lexer) stateSpace() stateFn {
  for unicode.IsSpace(l.peek()) {
    l.advance()
  }
  // We don't emit anything for spaces.
  return stateIdentifier
}

func (l *lexer) run(raw string) {
  for s := l.stateIdentifier(); s != nil {
    s = s()
  }
  close(l.items)
}

func (l *lexer) Lex(lval *yySymType) int {
  token, ok := <-l.next
  if !ok {
    return 0 // EOF
  }

  lval.raw = token

  log.Printf("%v", lval)

  switch token {
  case "(":
    return tokOpen
  case ")":
    return tokClose
  default:
    return tokIdentifier
  }
}

func (l *lexer) Error(s string) {
  log.Printf("parse error: %s\n", s)
}

func newLexer(raw string) *lexer {
  l := &lexer{
    next: make(chan item),
  }
  go l.run()
  return l
}

func Parse(raw string) {
  yyParse(newLexer(raw))
}
