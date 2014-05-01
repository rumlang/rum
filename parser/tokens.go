package parser

import (
	"fmt"

	"github.com/palats/glop/nodes"
)

const (
	_ = iota
	tokEOF
	tokOpen
	tokClose
	tokIdentifier
	tokInteger
	tokSpace
)

type tokenID int

func (t tokenID) String() string {
	switch t {
	case tokEOF:
		return "EOF"
	case tokOpen:
		return "Open"
	case tokClose:
		return "Close"
	case tokIdentifier:
		return "Identifier"
	case tokInteger:
		return "Integer"
	case tokSpace:
		return "Space"
	default:
		return fmt.Sprintf("Unknown[%d", t)
	}
}

var tokenPriorities = map[tokenID]int{
	tokEOF:  0,
	tokOpen: 30,
	// Do not use 0 for closing parenthesis - this way we can trigger an error
	// when we don't reach the end of the stream (e.g., "a)b") instead of
	// considering that the parsing is done.
	tokClose:      5,
	tokIdentifier: 20,
	tokInteger:    20,
}

// SourceRef contains information to trace code back to its implementation.
type SourceRef struct {
	// Line indicates the line in the file. Starts at 1.
	Line int
	// Column indicates the rune index (ignoring invalid sequences) in the line.
	Column int
}

// tokenInfo give details about a token the lexer extracted - including
// information about where it comes from.
type tokenInfo struct {
	// raw is the representation of the token, identical to the input. It also
	// includes invalid bytes (i.e., not matching a valid unicode point).
	raw string
	// text is the list of valid runes of this token.
	text []rune
	// ref contains information about where the token comes from.
	ref SourceRef
	// id is the lexer token ID, using tok* symbols.
	id tokenID
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
			ctx.Error(Error{
				Msg:  fmt.Sprintf("invalid token - expected ')', got: %q", string(t.text)),
				Code: ErrMissingClosingParenthesis,
				Ref:  t.ref,
			})
		} else {
			ctx.Advance()
		}
		return []nodes.Node{nodes.NewExpr(sublist)}
	// case tokClose: // Shoud never happen
	case tokIdentifier:
		return []nodes.Node{nodes.NewIdentifier(t)}
	case tokInteger:
		return []nodes.Node{nodes.NewInteger(t)}
	case tokEOF:
		// Needed for when an open parenthesis (or similar) is just before the end
		// of the input.
		return []nodes.Node{}
	}

	ctx.Error(Error{
		Msg:  fmt.Sprintf("unexpected %q (token type %s) at the beginning of an expression", t.text, t.id),
		Code: ErrInvalidNudToken,
		Ref:  t.ref,
	})
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
			ctx.Error(Error{
				Msg:  fmt.Sprintf("invalid token - expected ')', got: %q", string(t.text)),
				Code: ErrMissingClosingParenthesis,
				Ref:  t.ref,
			})
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

	ctx.Error(Error{
		Msg:  fmt.Sprintf("unexpected %q (token type %s) in an expression", t.text, t.id),
		Code: ErrInvalidLedToken,
		Ref:  t.ref,
	})
	return left
}

func (t tokenInfo) Lbp() int {
	return tokenPriorities[t.id]
}
