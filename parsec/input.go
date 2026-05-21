package parsec

import "fmt"

// Input is the interface that parser input streams must implement.
// All methods must be pure — implementations should be immutable value types.
type Input interface {
	// Head returns the next rune and true, or (0, false) at end of input.
	Head() (rune, bool)
	// Advance returns a new Input with the position moved one rune forward.
	Advance() Input
	// IsEOF reports whether the input is exhausted.
	IsEOF() bool
	// Pos returns the current position (used for error reporting).
	Pos() int
}

// stringInput is the built-in Input implementation for string (rune slice) input.
type stringInput struct {
	src []rune
	pos int
}

// NewInput returns an Input over the given string.
func NewInput(s string) Input {
	return stringInput{src: []rune(s)}
}

func (in stringInput) Head() (rune, bool) {
	if in.pos >= len(in.src) {
		return 0, false
	}
	return in.src[in.pos], true
}

func (in stringInput) Advance() Input {
	return stringInput{src: in.src, pos: in.pos + 1}
}

func (in stringInput) IsEOF() bool {
	return in.pos >= len(in.src)
}

func (in stringInput) Pos() int {
	return in.pos
}

// ParseError records position and message of a parse failure.
type ParseError struct {
	Pos     int
	Message string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error at %d: %s", e.Pos, e.Message)
}

func newError(in Input, msg string) error {
	return &ParseError{Pos: in.Pos(), Message: msg}
}
