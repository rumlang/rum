package parser

import (
	"fmt"

	"github.com/palats/glop/nodes"
)

type Context interface {
	Expression(rbp int) (interface{}, error)
	Advance() (Token, error)
}

type Token interface {
	// Nud is 'Null denotation'
	Nud(ctx Context) (interface{}, error)
	// Led is 'Left denotation'
	Led(ctx Context, left interface{}) (interface{}, error)
	// Lbp is 'Left Binding Priority'
	Lbp() int
}

type Lexer interface {
	Next() (Token, error)
}

type Parser struct {
	lex Lexer

	// token contains the next token to be encountered.
	token Token
}

// Advance get the next token from the lexer and return what was the current
// one.
func (p *Parser) Advance() (Token, error) {
	t := p.token
	next, err := p.lex.Next()
	p.token = next
	return t, err
}

func (p *Parser) Expression(rbp int) (interface{}, error) {
	t, err := p.Advance()
	if err != nil {
		return nil, err
	}
	left, err := t.Nud(p)
	if err != nil {
		return nil, err
	}
	for rbp < p.token.Lbp() {
		if t, err = p.Advance(); err != nil {
			return nil, err
		}
		if left, err = t.Led(p, left); err != nil {
			return nil, err
		}
	}
	return left, nil
}

func Parse(input string) (nodes.Node, error) {
	p := Parser{
		lex: newLexer(input),
	}
	// Make the first token available.
	p.Advance()
	n, err := p.Expression(0)
	if err != nil {
		return nil, err
	}
	result := n.([]nodes.Node)
	if len(result) == 0 {
		return nil, fmt.Errorf("no node found")
	}
	if len(result) != 1 {
		return nil, fmt.Errorf("obtained more than one node: %v", result)
	}
	return result[0], nil
}
