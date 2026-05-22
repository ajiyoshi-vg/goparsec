package parsec

// Parser consumes Input and returns a value, remaining Input, or an error.
type Parser[T any] func(Input) (T, Input, error)

// Run executes p on s, returning the parsed value (remaining input is ignored).
func Run[T any](p Parser[T], s string) (T, error) {
	v, rest, err := p(NewInput(s))
	if err == ErrNoMatch {
		// ErrNoMatch carries no message; rest holds the failure position.
		err = NewError(rest, "no match")
	}
	return v, err
}

// RunFull executes p on s and fails if any input remains after parsing.
func RunFull[T any](p Parser[T], s string) (T, error) {
	v, rest, err := p(NewInput(s))
	if err != nil {
		return v, err
	}
	if !rest.IsEOF() {
		return v, NewError(rest, "unexpected input")
	}
	return v, nil
}
