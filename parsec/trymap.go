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

// TryMap2 runs p1 and p2 in sequence, passing both results to f.
// If f returns an error, the parser fails with a ParseError.
func TryMap2[A, B, C any](p1 Parser[A], p2 Parser[B], f func(A, B) (C, error)) Parser[C] {
	return Bind(p1, func(a A) Parser[C] {
		return TryMap(p2, func(b B) (C, error) { return f(a, b) })
	})
}

// TryMap3 runs p1, p2, and p3 in sequence, passing all results to f.
// If f returns an error, the parser fails with a ParseError.
func TryMap3[A, B, C, D any](p1 Parser[A], p2 Parser[B], p3 Parser[C], f func(A, B, C) (D, error)) Parser[D] {
	return Bind(p1, func(a A) Parser[D] {
		return TryMap2(p2, p3, func(b B, c C) (D, error) { return f(a, b, c) })
	})
}

// TryMap4 runs p1, p2, p3, and p4 in sequence, passing all results to f.
// If f returns an error, the parser fails with a ParseError.
func TryMap4[A, B, C, D, E any](p1 Parser[A], p2 Parser[B], p3 Parser[C], p4 Parser[D], f func(A, B, C, D) (E, error)) Parser[E] {
	return Bind(p1, func(a A) Parser[E] {
		return TryMap3(p2, p3, p4, func(b B, c C, d D) (E, error) { return f(a, b, c, d) })
	})
}
