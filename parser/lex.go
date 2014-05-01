// Package parser takes care of lexing & parsing following the grammer defined
// in glop.y.
// This file contains the lexer, taking care of extracting relevant tokens from
// the provided input.
package parser

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/palats/glop/nodes"
)

// EOF is an arbitrary rune to indicate end of input from the lexer functions.
const EOF rune = 0

const (
	_ = iota
	tokEOF
	tokOpen
	tokClose
	tokIdentifier
	tokInteger
	tokSpace
)

var tokenPriorities = map[int]int{
	tokEOF:  0,
	tokOpen: 30,
	// Do not use 0 for closing parenthesis - this way we can trigger an error
	// when we don't reach the end of the stream (e.g., "a)b") instead of
	// considering that the parsing is done.
	tokClose:      5,
	tokIdentifier: 20,
	tokInteger:    20,
}

// stateFn is the prototype for function of the lexer state machine. They don't
// take any parameter - data is shared through the object.
type stateFn func() stateFn

// tokenInfo give details about a token the lexer extracted - including
// information about where it comes from.
type tokenInfo struct {
	// raw is the representation of the token, identical to the input. It also
	// includes invalid bytes (i.e., not matching a valid unicode point).
	raw string
	// text is the list of valid runes of this token.
	text []rune
	// startIndex is the absolute position in the input of this token, in bytes.
	startIndex int
	// startCol is the column this token start on the given line. This ignores
	// invalid bytes. First character on a line is column = 1.
	startCol int
	// line is the line number in the input of the beginning of this token.
	line int
	// id is the lexer token ID, using tok* symbols.
	id int
	// value is the parsed value of the token - can be a string, int, nil, ...
	value interface{}
}

// Value implements nodes.Token interface.
func (t tokenInfo) Value() interface{} {
	return t.value
}

// Nud implements the Token interface for the top down parser.
// It always returns a []nodes.Node{}. In case of errors, it will add an error
// through the context and return an empty list.
func (t tokenInfo) Nud(ctx Context) interface{} {
	switch t.id {
	case tokOpen:
		var sublist []nodes.Node
		if ctx.Peek().(tokenInfo).id != tokClose {
			sublist = ctx.Expression(tokenPriorities[tokClose]).([]nodes.Node)
		}
		t := ctx.Peek().(tokenInfo)
		if t.id != tokClose {
			ctx.Error(fmt.Errorf("invalid token - expected ')', got: %q", string(t.text)))
		} else {
			ctx.Advance()
		}
		return []nodes.Node{nodes.NewExpr(sublist)}
	// case tokClose: // Shoud never happen
	case tokIdentifier:
		return []nodes.Node{nodes.NewIdentifier(t)}
	case tokInteger:
		return []nodes.Node{nodes.NewInteger(t)}
	}

	ctx.Error(fmt.Errorf("unexpected %q (token id %d) at the beginning of an expression", t.text, t.id))
	return []nodes.Node{}
}

// Led implements the Token interface for the top down parser.
// It always returns a []nodes.Node{}. In case of errors, it will add an error
// through the context and ignore the token.
func (t tokenInfo) Led(ctx Context, left interface{}) interface{} {
	switch t.id {
	case tokOpen:
		var sublist []nodes.Node
		if ctx.Peek().(tokenInfo).id != tokClose {
			sublist = ctx.Expression(tokenPriorities[tokClose]).([]nodes.Node)
		}
		t := ctx.Peek().(tokenInfo)
		if t.id != tokClose {
			ctx.Error(fmt.Errorf("invalid token - expected ')', got: %q", string(t.text)))
		} else {
			ctx.Advance()
		}
		return append(left.([]nodes.Node), nodes.NewExpr(sublist))
	// case tokClose: // Should never happen.
	case tokIdentifier:
		return append(left.([]nodes.Node), nodes.NewIdentifier(t))
	case tokInteger:
		return append(left.([]nodes.Node), nodes.NewInteger(t))
	}

	ctx.Error(fmt.Errorf("unexpected %q (token id %d) in an expression", t.text, t.id))
	return left
}

func (t tokenInfo) Lbp() int {
	return tokenPriorities[t.id]
}

// lexer extract the tokens seen in the input.
type lexer struct {
	input string

	// next contains the next rune that would be added to current token with
	// avance().
	next rune
	// nextIndex is the byte index in the input of the 'next' rune.
	nextIndex int
	// nextCol is the rune index on the current line. Ignores invalid byte
	// sequences; 1-indexed.
	nextCol int
	// pos is the index in input of the next rune to be read (i.e., the index of
	// the rune starting after the rune currently in 'next')
	pos int
	// current line number of the rune in 'next'.
	line int

	// start indicates the index of the beginning of the token being parsed. This
	// is relative to input current data.
	start int
	// token is the token being built - it will be sent to the channel once
	// accept() is called.
	token   *tokenInfo
	tokens  chan tokenInfo
	errors  []string
	program nodes.Node
}

// peek looks one rune ahead in the input but does not advance the current
// pointer. It should never return RuneError.
func (l *lexer) peek() rune {
	return l.next
}

// advance moves the current position by one rune. Returns the rune encountered
// or utf8.RuneError if there was an issue. It should never return RuneError.
func (l *lexer) advance() rune {
	r := l.next
	l.token.text = append(l.token.text, r)
	if r == '\n' {
		l.line++
		l.nextCol = 1
	} else {
		l.nextCol++
	}

	invalid := ""
	found := false
	for !found && l.pos < len(l.input) {
		// w can be 0 when trying to decode an empty string. However, it should not
		// happen here because of the 'for' loop test.
		n, w := utf8.DecodeRuneInString(l.input[l.pos:])
		l.next = n
		l.nextIndex = l.pos
		l.pos += w
		if n == utf8.RuneError {
			// TODO: bail out if there are too many invalid bytes.
			invalid += l.input[l.pos : l.pos+w]
		} else {
			found = true
		}
	}

	if len(invalid) > 0 {
		// TODO: add errors with the invalid content.
	}

	if !found {
		l.next = EOF
	}
	return r
}

func (l *lexer) accept() tokenInfo {
	t := l.token
	l.token = &tokenInfo{
		startIndex: l.nextIndex,
		startCol:   l.nextCol,
		line:       l.line,
	}
	return *t
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
	if (len(token.text) > 1 && (token.text[0] == '+' || token.text[0] == '-')) || unicode.IsDigit(token.text[0]) {
		token.id = tokInteger
		i, err := strconv.ParseInt(string(token.text), 10, 64)
		if err != nil {
			panic(err) // TODO
		}
		token.value = i
	} else {
		token.id = tokIdentifier
		// Use []rune for identifiers?
		token.value = string(token.text)
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

func (l *lexer) Next() Token {
	token, ok := <-l.tokens
	if !ok {
		token = tokenInfo{
			id: tokEOF,
		}
	}
	return token
}

func newLexer(input string) *lexer {
	l := &lexer{
		input:  input,
		token:  &tokenInfo{},
		tokens: make(chan tokenInfo),
	}
	// Do an initial advance/accept to get the first character into 'next' and
	// make sure than the current token is properly initialized.
	l.advance()
	l.accept()
	go l.run()
	return l
}
