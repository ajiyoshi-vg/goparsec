package parsec

import (
	"errors"
	"fmt"

	"github.com/ajiyoshi-vg/goparsec/input"
)

// Positioned is implemented by errors that carry a source position.
// Choice uses Pos() to pick the furthest-reaching error among alternatives.
// User-defined error types may implement this interface to participate in
// Choice's furthest-error tracking, and to allow NewErrorAt to reconstruct
// a ParseError at that position without a *ParseError type assertion.
type Positioned interface {
	Pos() int  // rune offset from the start of input
	Line() int // 1-based line number
	Col() int  // 1-based column number
}

// ParseError records the position and message of a parse failure.
type ParseError struct {
	pos, line, col int
	Message        string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error at line %d, col %d: %s", e.line, e.col, e.Message)
}

// Pos, Line, Col implement Positioned.
func (e *ParseError) Pos() int  { return e.pos }
func (e *ParseError) Line() int { return e.line }
func (e *ParseError) Col() int  { return e.col }

// ErrNoMatch is the zero-allocation sentinel for a soft failure: the parser did not
// match at this position and consumed no input. The caller (Choice, Label) is
// responsible for converting it to a *ParseError when a user-visible message is needed.
//
// Custom parsers should return ErrNoMatch when they simply do not match the input
// at the current position, allowing Choice to try the next alternative without
// allocating an error object.
var ErrNoMatch = errors.New("parsec: no match")

// NewError creates a *ParseError at the current position of in with the given message.
// Custom parsers should use NewError to produce positioned errors that integrate
// correctly with Choice's furthest-error tracking and Label's message replacement.
func NewError(in input.Input, msg string) error {
	return &ParseError{pos: in.Pos(), line: in.Line(), col: in.Col(), Message: msg}
}

// NewErrorf is like NewError but formats the message using fmt.Sprintf.
func NewErrorf(in input.Input, format string, args ...any) error {
	return &ParseError{pos: in.Pos(), line: in.Line(), col: in.Col(), Message: fmt.Sprintf(format, args...)}
}

// NewErrorAt creates a *ParseError at the same position as p, with the given message.
// Use this inside MapError callbacks to replace an error message while preserving
// position — works with any error that implements Positioned, not just *ParseError.
func NewErrorAt(p Positioned, msg string) error {
	return &ParseError{pos: p.Pos(), line: p.Line(), col: p.Col(), Message: msg}
}
