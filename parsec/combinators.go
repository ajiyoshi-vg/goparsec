package parsec

import (
	"fmt"
	"strconv"
)

// Many parses zero or more occurrences of p. Always succeeds.
// Panics if p succeeds without consuming input, which would cause an infinite loop.
func Many[T any](p Parser[T]) Parser[[]T] {
	return func(in Input) ([]T, Input, error) {
		results := []T{}
		cur := in
		for {
			val, next, err := p(cur)
			if err != nil {
				return results, cur, nil
			}
			if next.Pos() == cur.Pos() {
				panic("parsec: Many: parser succeeded without consuming input")
			}
			results = append(results, val)
			cur = next
		}
	}
}

// Many1 parses one or more occurrences of p.
func Many1[T any](p Parser[T]) Parser[[]T] {
	return func(in Input) ([]T, Input, error) {
		val, cur, err := p(in)
		if err != nil {
			return nil, in, err
		}
		rest, cur, _ := Many(p)(cur)
		return append([]T{val}, rest...), cur, nil
	}
}

// Choice tries each parser in order, returning the first success.
// On failure, no input is consumed (all parsers implicitly backtrack).
// Reports the error from whichever parser reached the furthest position.
func Choice[T any](ps ...Parser[T]) Parser[T] {
	return func(in Input) (T, Input, error) {
		var zero T
		var bestErr error
		for _, p := range ps {
			val, next, err := p(in)
			if err == nil {
				return val, next, nil
			}
			bestErr = furthestError(bestErr, err)
		}
		return zero, in, bestErr
	}
}

// furthestError returns whichever error occurred at the greater input position.
func furthestError(a, b error) error {
	if a == nil {
		return b
	}
	pa, ok1 := a.(*ParseError)
	pb, ok2 := b.(*ParseError)
	if !ok1 || !ok2 {
		return b
	}
	if pb.Pos >= pa.Pos {
		return b
	}
	return a
}

// Return always succeeds with v without consuming any input.
func Return[T any](v T) Parser[T] {
	return func(in Input) (T, Input, error) {
		return v, in, nil
	}
}

// Option returns def if p fails, without consuming input.
func Option[T any](def T, p Parser[T]) Parser[T] {
	return func(in Input) (T, Input, error) {
		val, next, err := p(in)
		if err != nil {
			return def, in, nil
		}
		return val, next, nil
	}
}

// Map transforms the result of p using f.
func Map[T, U any](p Parser[T], f func(T) U) Parser[U] {
	return func(in Input) (U, Input, error) {
		val, next, err := p(in)
		if err != nil {
			var zero U
			return zero, in, err
		}
		return f(val), next, nil
	}
}

// Bind sequences p then applies f to its result to obtain the next parser.
func Bind[T, U any](p Parser[T], f func(T) Parser[U]) Parser[U] {
	return func(in Input) (U, Input, error) {
		val, next, err := p(in)
		if err != nil {
			var zero U
			return zero, in, err
		}
		return f(val)(next)
	}
}

// Then runs pa, discards its result, then runs pb.
func Then[T, U any](pa Parser[T], pb Parser[U]) Parser[U] {
	return func(in Input) (U, Input, error) {
		_, next, err := pa(in)
		if err != nil {
			var zero U
			return zero, in, err
		}
		return pb(next)
	}
}

// Skip runs pa then pb, returning pa's result.
func Skip[T, U any](pa Parser[T], pb Parser[U]) Parser[T] {
	return func(in Input) (T, Input, error) {
		val, next, err := pa(in)
		if err != nil {
			return val, in, err
		}
		_, next, err = pb(next)
		if err != nil {
			var zero T
			return zero, in, err
		}
		return val, next, nil
	}
}

// Between parses open, then p, then close, returning p's result.
func Between[O, C, T any](open Parser[O], close Parser[C], p Parser[T]) Parser[T] {
	return Then(open, Skip(p, close))
}

// SepBy parses zero or more occurrences of p separated by sep.
func SepBy[T, S any](p Parser[T], sep Parser[S]) Parser[[]T] {
	return func(in Input) ([]T, Input, error) {
		val, cur, err := p(in)
		if err != nil {
			return []T{}, in, nil
		}
		results := []T{val}
		for {
			_, next, err := sep(cur)
			if err != nil {
				return results, cur, nil
			}
			val, next, err = p(next)
			if err != nil {
				return results, cur, nil
			}
			results = append(results, val)
			cur = next
		}
	}
}

// Spaces parses zero or more whitespace characters.
func Spaces() Parser[string] {
	return Map(Many(Space()), func(rs []rune) string { return string(rs) })
}

// Chainl1 parses one or more occurrences of p separated by op,
// folding the results left-associatively.
func Chainl1[T any](p Parser[T], op Parser[func(T, T) T]) Parser[T] {
	type step struct {
		fn  func(T, T) T
		val T
	}
	steps := Many(Bind(op, func(fn func(T, T) T) Parser[step] {
		return Map(p, func(v T) step { return step{fn, v} })
	}))
	return Bind(p, func(acc T) Parser[T] {
		return Map(steps, func(ss []step) T {
			for _, s := range ss {
				acc = s.fn(acc, s.val)
			}
			return acc
		})
	})
}

// Integer parses an optional '-' followed by one or more digits.
func Integer() Parser[int] {
	neg := Option(false, Map(Char('-'), func(rune) bool { return true }))
	return Bind(neg, func(isNeg bool) Parser[int] {
		return Map(Natural(), func(n int) int {
			if isNeg {
				return -n
			}
			return n
		})
	})
}

// Float parses an optional '-', one or more digits, an optional fractional part
// ('.' digits), and an optional exponent ([eE][+-]? digits), returning float64.
func Float() Parser[float64] {
	sign := Option("", Map(Char('-'), func(rune) string { return "-" }))
	digits := Map(Many1(Digit()), func(rs []rune) string { return string(rs) })
	frac := Option("", Map(
		Bind(Char('.'), func(rune) Parser[[]rune] { return Many1(Digit()) }),
		func(ds []rune) string { return "." + string(ds) },
	))
	expSign := Option("", Choice(
		Map(Char('+'), func(rune) string { return "+" }),
		Map(Char('-'), func(rune) string { return "-" }),
	))
	exp := Option("", Bind(
		Satisfy(func(r rune) bool { return r == 'e' || r == 'E' }),
		func(e rune) Parser[string] {
			return Bind(expSign, func(s string) Parser[string] {
				return Map(digits, func(d string) string { return string(e) + s + d })
			})
		},
	))
	floatStr := Bind(sign, func(s string) Parser[string] {
		return Bind(digits, func(d string) Parser[string] {
			return Bind(frac, func(f string) Parser[string] {
				return Map(exp, func(e string) string { return s + d + f + e })
			})
		})
	})
	return Map(floatStr, func(s string) float64 {
		f, _ := strconv.ParseFloat(s, 64)
		return f
	})
}

// NotFollowedBy succeeds if p fails, consuming no input.
// Useful for ensuring a token is not followed by unexpected characters
// (e.g., distinguishing keywords from identifiers).
func NotFollowedBy[T any](p Parser[T]) Parser[struct{}] {
	return func(in Input) (struct{}, Input, error) {
		_, _, err := p(in)
		if err == nil {
			c, ok := in.Head()
			if ok {
				return struct{}{}, in, newError(in, fmt.Sprintf("unexpected %q", c))
			}
			return struct{}{}, in, newError(in, "unexpected input")
		}
		return struct{}{}, in, nil
	}
}

// Label attaches a description to p, replacing its error message on failure.
func Label[T any](p Parser[T], label string) Parser[T] {
	return func(in Input) (T, Input, error) {
		val, next, err := p(in)
		if err != nil {
			return val, in, &ParseError{Pos: in.Pos(), Message: "expected " + label}
		}
		return val, next, nil
	}
}

// Count parses exactly n occurrences of p, failing if fewer are found.
func Count[T any](n int, p Parser[T]) Parser[[]T] {
	return func(in Input) ([]T, Input, error) {
		results := make([]T, 0, n)
		cur := in
		for range n {
			val, next, err := p(cur)
			if err != nil {
				return nil, in, err
			}
			results = append(results, val)
			cur = next
		}
		return results, cur, nil
	}
}

// ManyTill parses p zero or more times until end succeeds, consuming end.
func ManyTill[T, E any](p Parser[T], end Parser[E]) Parser[[]T] {
	return func(in Input) ([]T, Input, error) {
		var results []T
		cur := in
		for {
			if _, next, err := end(cur); err == nil {
				return results, next, nil
			}
			val, next, err := p(cur)
			if err != nil {
				return nil, in, err
			}
			results = append(results, val)
			cur = next
		}
	}
}

// Chainr1 parses one or more occurrences of p separated by op,
// folding the results right-associatively.
func Chainr1[T any](p Parser[T], op Parser[func(T, T) T]) Parser[T] {
	var rec Parser[T]
	rec = func(in Input) (T, Input, error) {
		return Bind(p, func(x T) Parser[T] {
			return Option(x, Bind(op, func(f func(T, T) T) Parser[T] {
				return Map(rec, func(y T) T { return f(x, y) })
			}))
		})(in)
	}
	return rec
}

// Natural parses one or more digits as a non-negative integer.
func Natural() Parser[int] {
	return Map(Many1(Digit()), func(digits []rune) int {
		n := 0
		for _, d := range digits {
			n = n*10 + int(d-'0')
		}
		return n
	})
}
