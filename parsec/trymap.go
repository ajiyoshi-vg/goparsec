package parsec

import "github.com/ajiyoshi-vg/goparsec/input"

// TryMap runs p, then passes the result to f. If f returns an error,
// the parser fails with a ParseError at the position where p ended.
func TryMap[T, U any](p Parser[T], f func(T) (U, error)) Parser[U] {
	return func(in input.Input) (U, input.Input, error) {
		v, rest, err := p(in)
		if err != nil {
			var zero U
			return zero, in, err
		}
		u, err := f(v)
		if err != nil {
			var zero U
			return zero, in, NewError(rest, err.Error())
		}
		return u, rest, nil
	}
}
