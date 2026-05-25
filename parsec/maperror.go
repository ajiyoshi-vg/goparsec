package parsec

import "github.com/ajiyoshi-vg/goparsec/input"

// MapError runs p and, on failure, calls f with the entry position and the
// error to produce a replacement error. On success, the result passes through
// unchanged. f receives the input position at the start of p (not where p
// stopped), so it can create a new positioned error or inspect and forward the
// original.
func MapError[T any](p Parser[T], f func(input.Input, error) error) Parser[T] {
	return func(in input.Input) (T, input.Input, error) {
		val, next, err := p(in)
		if err == nil {
			return val, next, nil
		}
		return val, in, f(in, err)
	}
}
