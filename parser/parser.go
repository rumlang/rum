package parser

import (
	"fmt"

	"github.com/palats/glop/nodes"
)

type Context interface {
	Expression(rbp int) interface{}
	Advance() Token
}

type Token interface {
	// Nud is 'Null denotation'
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
}

// Advance get the next token from the lexer and return what was the current
// one.
func (p *Parser) Advance() Token {
	t := p.token
	p.token = p.lex.Next()
	return t
}

func (p *Parser) Expression(rbp int) interface{} {
	t := p.Advance()
	left := t.Nud(p)
	for rbp < p.token.Lbp() {
		t = p.Advance()
		left = t.Led(p, left)
	}
	return left
}

func Parse(input string) (nodes.Node, error) {
	p := Parser{
		lex: newLexer(input),
	}
	// Make the first token available.
	p.Advance()
	result := p.Expression(0).([]nodes.Node)
	if len(result) == 0 {
		return nil, fmt.Errorf("no node found")
	}
	if len(result) != 1 {
		return nil, fmt.Errorf("obtained more than one node: %v", result)
	}
	return result[0], nil
}
