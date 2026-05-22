package parsec

// Parser consumes Input and returns a value, remaining Input, or an error.
type Parser[T any] func(Input) (T, Input, error)

// Run executes p on in, returning the parsed value (remaining input is ignored).
func Run[T any](p Parser[T], in Input) (T, error) {
	v, rest, err := p(in)
	if err == ErrNoMatch {
		err = NewError(rest, "no match")
	}
	return v, err
}

// RunFull executes p on in and fails if any input remains after parsing.
func RunFull[T any](p Parser[T], in Input) (T, error) {
	v, rest, err := p(in)
	if err != nil {
		return v, err
	}
	if !rest.IsEOF() {
		return v, NewError(rest, "unexpected input")
	}
	return v, nil
}

// RunString executes p on s, returning the parsed value (remaining input is ignored).
// Shorthand for Run(p, NewStringInput(s)).
func RunString[T any](p Parser[T], s string) (T, error) {
	return Run(p, NewStringInput(s))
}

// RunStringFull executes p on s and fails if any input remains after parsing.
// Shorthand for RunFull(p, NewStringInput(s)).
func RunStringFull[T any](p Parser[T], s string) (T, error) {
	return RunFull(p, NewStringInput(s))
}
