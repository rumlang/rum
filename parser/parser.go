package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/palats/glop/nodes"
)

func Parse(input string) (nodes.Node, error) {
	p := NewTopDown(newLexer(input))
	result := p.Expression(tokenPriorities[tokEOF]).([]nodes.Node)
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
