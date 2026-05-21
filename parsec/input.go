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
	// Pos returns the rune offset from the start (used for furthest-error comparison).
	Pos() int
	// Line returns the 1-based line number of the current position.
	Line() int
	// Col returns the 1-based column number of the current position.
	Col() int
}

// stringInput is the built-in Input implementation for string (rune slice) input.
type stringInput struct {
	src  []rune
	pos  int
	line int
	col  int
}

// NewInput returns an Input over the given string.
func NewInput(s string) Input {
	return stringInput{src: []rune(s), line: 1, col: 1}
}

func (in stringInput) Head() (rune, bool) {
	if in.pos >= len(in.src) {
		return 0, false
	}
	return in.src[in.pos], true
}

func (in stringInput) Advance() Input {
	if in.pos >= len(in.src) {
		return in
	}
	next := stringInput{src: in.src, pos: in.pos + 1}
	if in.src[in.pos] == '\n' {
		next.line = in.line + 1
		next.col = 1
	} else {
		next.line = in.line
		next.col = in.col + 1
	}
	return next
}

func (in stringInput) IsEOF() bool { return in.pos >= len(in.src) }
func (in stringInput) Pos() int    { return in.pos }
func (in stringInput) Line() int   { return in.line }
func (in stringInput) Col() int    { return in.col }

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

func newError(in Input, msg string) error {
	return &ParseError{Pos: in.Pos(), Line: in.Line(), Col: in.Col(), Message: msg}
}
