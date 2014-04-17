package parser

import (
	"errors"
	"fmt"

	"github.com/palats/glop/nodes"
)

func Parse(input string) (nodes.Node, []error) {
	r, errs := TopDownParse(newLexer(input))
	result := r.([]nodes.Node)
	var n nodes.Node
	if len(result) == 0 {
		errs = append(errs, errors.New("no node found"))
	} else if len(result) != 1 {
		errs = append(errs, fmt.Errorf("obtained more than one node: %v", result))
	} else {
		n = result[0]
	}

	return n, errs
}
