package parsec

import "github.com/ajiyoshi-vg/goparsec/input"

// Expect runs p and replaces any failure's message with label.
// Unlike Label, Expect also replaces hard failure messages.
// The error position is preserved: soft failures report at the entry
// position; hard failures report at the position where p stopped.
//
// Nesting behaviour:
//   - Label(Label(p, "inner"), "outer") — inner wins (outer cannot replace hard errors)
//   - Expect(Expect(p, "inner"), "outer") — outer wins (always replaces)
func Expect[T any](p Parser[T], label string) Parser[T] {
	return func(in input.Input) (T, input.Input, error) {
		val, next, err := p(in)
		if err == nil {
			return val, next, nil
		}
		if err == ErrNoMatch {
			return val, in, NewError(in, label)
		}
		// Hard failure: replace message, preserve position.
		if pe, ok := err.(*ParseError); ok {
			return val, in, &ParseError{
				Pos:     pe.Pos,
				Line:    pe.Line,
				Col:     pe.Col,
				Message: label,
			}
		}
		return val, in, NewError(in, label)
	}
}
