package parsec

import (
	"errors"
	"fmt"

	"github.com/ajiyoshi-vg/goparsec/input"
)

// ParseError records the position and message of a parse failure.
type ParseError struct {
	Pos     int // rune offset, used for furthest-error comparison
	Line    int // 1-based line number
	Col     int // 1-based column number
	Message string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error at line %d, col %d: %s", e.Line, e.Col, e.Message)
}

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
	return &ParseError{Pos: in.Pos(), Line: in.Line(), Col: in.Col(), Message: msg}
}

// NewErrorf is like NewError but formats the message using fmt.Sprintf.
func NewErrorf(in input.Input, format string, args ...any) error {
	return &ParseError{Pos: in.Pos(), Line: in.Line(), Col: in.Col(), Message: fmt.Sprintf(format, args...)}
}
