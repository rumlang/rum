package parser

import (
	"log"
	"strconv"
	"unicode"
	"unicode/utf8"

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
	text  string
	value interface{}
	id    int
}

func (t tokenInfo) Value() interface{} {
	return t.value
}

// lexer is a 'go yacc' compatible lexer object.
type lexer struct {
	input   string
	start   int
	pos     int
	tokens  chan tokenInfo
	errors  []string
	program nodes.Node
}

func (l *lexer) peek() rune {
	if l.pos >= len(l.input) {
		return EOF
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.pos:])
	// TODO error detection
	return r
}

func (l *lexer) advance() {
	_, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += w
}

func (l *lexer) accept() tokenInfo {
	t := tokenInfo{
		text: l.input[l.start:l.pos],
	}
	l.start = l.pos
	return t
}

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
	log.Printf("%v", lval)
	return token.id
}

func (l *lexer) Error(s string) {
	l.errors = append(l.errors, s)
	log.Printf("parse error: %s\n", s)
}

func newLexer(input string) *lexer {
	l := &lexer{
		input:  input,
		tokens: make(chan tokenInfo),
	}
	go l.run()
	return l
}

func Parse(input string) nodes.Node {
	l := newLexer(input)
	yyParse(l)
	return l.program
}
