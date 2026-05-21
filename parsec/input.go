package parsec

import "fmt"

// Input is an immutable view into the parser input stream.
type Input struct {
	src []rune
	pos int
}

func NewInput(s string) Input {
	return Input{src: []rune(s)}
}

func (in Input) Head() (rune, bool) {
	if in.pos >= len(in.src) {
		return 0, false
	}
	return in.src[in.pos], true
}

func (in Input) Advance() Input {
	return Input{src: in.src, pos: in.pos + 1}
}

func (in Input) IsEOF() bool {
	return in.pos >= len(in.src)
}

func (in Input) Pos() int {
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
	return &ParseError{Pos: in.pos, Message: msg}
}
