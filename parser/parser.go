package parser

import (
	"errors"
	"fmt"
	"unicode/utf8"

	"github.com/palats/glop/nodes"
)

// Source holds a rune representation of source code.
type Source struct {
	// data contains all the runes found in the source.
	data []rune
	// lines are the index within data of the beginning of each line.
	lines []int
}

// Scan returns a channel providing all the rune of the source in order one at
// a time. Channel will be closed at the end of the source.
func (s *Source) Scan() <-chan rune {
	ch := make(chan rune)
	go func() {
		for _, c := range s.data {
			ch <- c
		}
		close(ch)
	}()
	return ch
}

// Line returns the line corresponding to the 0-based provided index. Returns a
// non-nil error if the value is out of bound.
func (s *Source) Line(i int) ([]rune, error) {
	if i < 0 || i >= len(s.lines) {
		return nil, fmt.Errorf("out of bound line %d (source has %d lines)", i, len(s.lines))
	}

	begin := s.lines[i]
	var end int
	if i == len(s.lines)-1 {
		// Last line
		end = len(s.data)
	} else {
		// Other lines
		end = s.lines[i+1]
	}
	return s.data[begin:end], nil
}

// NewSource creates a new source object from the provided utf8 string. It will
// ignore all invalid codepoint.
func NewSource(input string) *Source {
	s := &Source{
		lines: []int{0},
	}

	for _, c := range input {
		if c == utf8.RuneError {
			// TODO
			continue
		}
		s.data = append(s.data, c)
		if c == '\n' {
			s.lines = append(s.lines, len(s.data))
		}
	}

	return s
}

// SourceRef contains information to trace code back to its implementation.
type SourceRef struct {
	// Source to which this refers to.
	Source *Source
	// Line indicates the line in the file. 0-indexed.
	Line int
	// Column indicates the rune index (ignoring invalid sequences) in the line.
	// 0-indexed.
	Column int
}

// Parse will take the provided source, parse it, and ensure that only one root
// node is returned.
func Parse(src *Source) (nodes.Node, []error) {
	r, errs := TopDownParse(newLexer(src))
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
