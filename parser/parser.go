package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/palats/glop/nodes"
)

type Context interface {
	Expression(rbp int) interface{}
	Advance() Token
	Peek() Token
	Error(string)
}

type Token interface {
	// Nud is 'Null denotation'. If it returns nil, the parser will skip it and
	// use the next token to start the expression.
	Nud(ctx Context) interface{}
	// Led is 'Left denotation'
	Led(ctx Context, left interface{}) interface{}
	// Lbp is 'Left Binding Priority'
	Lbp() int
}

type Lexer interface {
	Next() Token
}

type Parser struct {
	lex Lexer

	// token contains the next token to be encountered.
	token Token

	// All errors coming from the tokens nud/led functions.
	errors []string
}

// Advance get the next token from the lexer and return what was the current
// one.
func (p *Parser) Advance() Token {
	t := p.token
	p.token = p.lex.Next()
	return t
}

func (p *Parser) Peek() Token {
	return p.token
}

func (p *Parser) Expression(rbp int) interface{} {
	var left interface{}

	for left == nil {
		left = p.Advance().Nud(p)
	}
	for rbp < p.Peek().Lbp() {
		left = p.Advance().Led(p, left)
	}
	return left
}

func (p *Parser) Error(s string) {
	p.errors = append(p.errors, s)
}

func Parse(input string) (nodes.Node, error) {
	p := Parser{
		lex: newLexer(input),
	}
	// Make the first token available.
	p.Advance()
	result := p.Expression(0).([]nodes.Node)
	var n nodes.Node
	if len(result) == 0 {
		p.Error("no node found")
	} else if len(result) != 1 {
		p.Error(fmt.Sprintf("obtained more than one node: %v", result))
	} else {
		n = result[0]
	}

	var err error
	if len(p.errors) > 0 {
		err = errors.New(strings.Join(p.errors, " -- "))
	}
	return n, err
}
