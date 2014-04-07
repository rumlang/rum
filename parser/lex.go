// Package parser takes care of lexing & parsing following the grammer defined
// in glop.y.
// This file contains the lexer, taking care of extracting relevant tokens from
// the provided input.
package parser

import (
	"strconv"
	"unicode"
	"unicode/utf8"

	log "github.com/golang/glog"
	"github.com/palats/glop/nodes"
)

// EOF is an arbitrary rune to indicate end of input from the lexer functions.
const EOF rune = 0

// stateFn is the prototype for function of the lexer state machine. They don't
// take any parameter - data is shared through the object.
type stateFn func() stateFn

// tokenInfo give details about a token the lexer extracted - including
// information about where it comes from.
type tokenInfo struct {
	// text is the raw representation, identical to the input.
	text string
	// start is the absolute position in the input of this token.
	start int
	// line is the line number in the input of the beginning of this token.
	line int
	// value is the parsed value of the token - can be string, int, nil, ...
	value interface{}
	// id is the lexer token ID, using tok* symbols defined in glop.y.
	id int
}

// Value implements nodes.Token interface.
func (t tokenInfo) Value() interface{} {
	return t.value
}

// lexer is a 'go yacc' compatible lexer object.
type lexer struct {
	input string

	// start indicates the index of the beginning of the token being parsed. This
	// is relative to input current data.
	start int
	// pos is the index in input of the next rune to be read.
	pos int
	// current line number of the start position.
	line    int
	tokens  chan tokenInfo
	errors  []string
	program nodes.Node
}

// peek looks one rune ahead in the input but does not advance the current
// pointer. If the input is invalid, it will return utf8.RuneError.
func (l *lexer) peek() rune {
	if l.pos >= len(l.input) {
		return EOF
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.pos:])
	// TODO error detection
	return r
}

// advance moves the current position by one rune. Returns the rune encountered
// or utf8.RuneError if there was an issue.
func (l *lexer) advance() rune {
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += w
	return r
}

func (l *lexer) accept() tokenInfo {
	t := tokenInfo{
		text:  l.input[l.start:l.pos],
		start: l.start,
		line:  l.line,
	}
	l.start = l.pos

	for _, r := range t.text {
		if r == '\n' {
			l.line++
		}
	}

	return t
}

// stateIdentifier parses arbitrary strings & numbers.
func (l *lexer) stateIdentifier() stateFn {
	var next stateFn
	for next == nil {
		r := l.peek()

		switch {
		case r == '(':
			next = l.stateOpen
		case r == ')':
			next = l.stateClose
		case unicode.IsSpace(r):
			next = l.stateSpace
		case r == EOF:
			next = l.stateEnd
		case r == utf8.RuneError:
			//
		default:
			l.advance()
		}
	}

	token := l.accept()

	// Ignore empty transition - they're just a parsing artifact.
	if len(token.text) == 0 {
		return next
	}

	// Check the first rune to determine whether it is just an arbitrary
	// identifier or a number.
	r, w := utf8.DecodeRuneInString(token.text)

	if (len(token.text)-w > 0 && (r == '+' || r == '-')) || unicode.IsDigit(r) {
		token.id = tokInteger
		i, err := strconv.ParseInt(token.text, 10, 64)
		if err != nil {
			panic(err) // TODO
		}
		token.value = i
	} else {
		token.id = tokIdentifier
		token.value = token.text
	}
	l.tokens <- token
	return next
}

func (l *lexer) stateOpen() stateFn {
	l.advance()
	token := l.accept()
	// TODO: check that it is the right character and fail otherwise.
	token.id = tokOpen
	l.tokens <- token
	return l.stateIdentifier
}

func (l *lexer) stateClose() stateFn {
	l.advance()
	token := l.accept()
	// TODO: check that it is the right character and fail otherwise.
	token.id = tokClose
	l.tokens <- token
	return l.stateIdentifier
}

func (l *lexer) stateSpace() stateFn {
	for unicode.IsSpace(l.peek()) {
		l.advance()
	}
	// We don't emit anything for spaces.
	l.accept()
	return l.stateIdentifier
}

func (l *lexer) stateEnd() stateFn {
	return nil
}

func (l *lexer) run() {
	for s := l.stateIdentifier(); s != nil; {
		s = s()
	}
	close(l.tokens)
}

func (l *lexer) Lex(lval *yySymType) int {
	token, ok := <-l.tokens
	if !ok {
		return 0 // EOF
	}

	lval.token = token
	return token.id
}

func (l *lexer) Error(s string) {
	l.errors = append(l.errors, s)
	log.Errorf("parse error: %s\n", s)
}

func newLexer(input string) *lexer {
	l := &lexer{
		input:  input,
		tokens: make(chan tokenInfo),
	}
	go l.run()
	return l
}

// Parse analyzes the provided input and returns the correspond AST.
func Parse(input string) nodes.Node {
	l := newLexer(input)
	yyParse(l)
	return l.program
}
