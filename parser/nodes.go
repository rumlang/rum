// Package parser implements the glop parser.
// This file implements the nodes of the AST.
package parser

import (
	"fmt"
	"strings"
)

const (
	// NodeExpression - Value is []*Node
	NodeExpression = iota
	// NodeIdentifier - Value is string
	NodeIdentifier
	// NodeInteger - Value is int64
	NodeInteger
	// NodeFloat - Value is float64
	NodeFloat
)

type NodeType int

func (t NodeType) String() string {
	switch t {
	case NodeExpression:
		return "Expression"
	case NodeIdentifier:
		return "Identifier"
	case NodeInteger:
		return "Integer"
	default:
		return fmt.Sprintf("Unknown[%d]", t)
	}
}

// Node provides information about one node of the AST.
type Node struct {
	Type  NodeType
	Value interface{}
	Ref   SourceRef
}

// Children returns the list of children this Node has. Most node cannot have
// any and so return nil.
func (n *Node) Children() []*Node {
	if n.Type != NodeExpression {
		return nil
	}
	return n.Value.([]*Node)
}

// String provides a type dependent representation of the node.
func (n *Node) String() string {
	switch n.Type {
	case NodeExpression:
		var elt []string
		for _, node := range n.Children() {
			elt = append(elt, node.String())
		}
		return fmt.Sprintf("<expr>(%s)", strings.Join(elt, " "))
	case NodeIdentifier:
		return fmt.Sprintf("<id>%q", n.Value.(string))
	case NodeInteger:
		return fmt.Sprintf("<integer>%d", n.Value.(int64))
	case NodeFloat:
		return fmt.Sprintf("<float>%f", n.Value.(float64))
	default:
		panic(fmt.Sprintf("Forgot to support a new type of node it seems: %v", n.Type))
	}
}

func NewNode(t NodeType, v interface{}, ref SourceRef) *Node {
	return &Node{
		Type:  t,
		Value: v,
		Ref:   ref,
	}
}
