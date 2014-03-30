package parser

import (
  "log"
  "unicode"
  "unicode/utf8"
)

const EOF rune = 0

type stateFn func() stateFn

type tokenInfo struct {
  raw string
  id int
}

type lexer struct {
  input string
  start int
  pos int
  tokens chan tokenInfo
  errors []string
}

func (l *lexer) peek() rune {
  if l.pos >= len(l.input) {
    return EOF
  }
  r, _ := utf8.DecodeRuneInString(l.input[l.pos:])
  // XXX error detection
  return r
}

func (l *lexer) advance() {
  _, w := utf8.DecodeRuneInString(l.input[l.pos:])
  l.pos += w
}

func (l *lexer) emit(token int) {
  if l.pos <= l.start {
    return
  }
  s := l.input[l.start:l.pos]
  l.tokens <- tokenInfo{raw: s, id: token}
  l.start = l.pos
}

func (l *lexer) discard() {
  l.start = l.pos
}

func (l *lexer) stateIdentifier() stateFn {
  for {
    r := l.peek()

    switch {
    case r == '(':
      l.emit(tokIdentifier)
      l.advance()
      l.emit(tokOpen)
    case r == ')':
      l.emit(tokIdentifier)
      l.advance()
      l.emit(tokClose)
    case unicode.IsSpace(r):
      l.emit(tokIdentifier)
      return l.stateSpace
    case r == EOF:
      l.emit(tokIdentifier)
      return nil
    default:
      l.advance()
    }
  }
}

func (l *lexer) stateSpace() stateFn {
  for unicode.IsSpace(l.peek()) {
    l.advance()
  }
  // We don't emit anything for spaces.
  l.discard()
  return l.stateIdentifier
}

func (l *lexer) run() {
  for s := l.stateIdentifier(); s != nil ; {
    s = s()
  }
  close(l.tokens)
}

func (l *lexer) Lex(lval *yySymType) int {
  token, ok := <-l.tokens
  if !ok {
    return 0 // EOF
  }

  lval.raw = token.raw
  log.Printf("%v", lval)
  return token.id
}

func (l *lexer) Error(s string) {
  l.errors = append(l.errors, s)
  log.Printf("parse error: %s\n", s)
}

func newLexer(raw string) *lexer {
  l := &lexer{
    input: raw,
    tokens: make(chan tokenInfo),
  }
  go l.run()
  return l
}

func Parse(raw string) {
  yyParse(newLexer(raw))
}
