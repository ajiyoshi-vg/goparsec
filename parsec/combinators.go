package parsec

import (
	"fmt"
	"strconv"
	"strings"
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
		results := []T{val}
		for {
			v, next, err := p(cur)
			if err != nil {
				return results, cur, nil
			}
			if next.Pos() == cur.Pos() {
				panic("parsec: Many1: parser succeeded without consuming input")
			}
			results = append(results, v)
			cur = next
		}
	}
}

// Choice tries each parser in order, returning the first success.
// On failure, no input is consumed (all parsers implicitly backtrack).
// Reports the error from whichever parser reached the furthest position.
// Fails immediately if called with no parsers.
func Choice[T any](ps ...Parser[T]) Parser[T] {
	return func(in Input) (T, Input, error) {
		var zero T
		if len(ps) == 0 {
			return zero, in, NewError(in, "Choice: no alternatives")
		}
		var bestErr error
		for _, p := range ps {
			val, next, err := p(in)
			if err == nil {
				return val, next, nil
			}
			bestErr = furthestError(bestErr, err)
		}
		// All alternatives were soft failures with no position info: synthesize a *ParseError.
		if bestErr == ErrNoMatch {
			return zero, in, NewError(in, "no alternatives matched")
		}
		return zero, in, bestErr
	}
}

// furthestError returns whichever error occurred at the greater input position.
// ErrNoMatch (soft failure with no position) always loses to a *ParseError.
func furthestError(a, b error) error {
	if a == nil || a == ErrNoMatch {
		return b
	}
	if b == nil || b == ErrNoMatch {
		return a
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
	return func(in Input) (int, Input, error) {
		cur := in
		neg := false
		if c, ok := cur.Head(); ok && c == '-' {
			neg = true
			cur = cur.Advance()
		}
		c, ok := cur.Head()
		if !ok || c < '0' || c > '9' {
			return 0, in, ErrNoMatch
		}
		n := int(c - '0')
		cur = cur.Advance()
		for {
			c, ok = cur.Head()
			if !ok || c < '0' || c > '9' {
				break
			}
			n, cur = n*10+int(c-'0'), cur.Advance()
		}
		if neg {
			return -n, cur, nil
		}
		return n, cur, nil
	}
}

// Float parses an optional '-', one or more digits, an optional fractional part
// ('.' digits), and an optional exponent ([eE][+-]? digits), returning float64.
func Float() Parser[float64] {
	return func(in Input) (float64, Input, error) {
		var b strings.Builder
		cur := in

		if c, ok := cur.Head(); ok && c == '-' {
			b.WriteByte('-')
			cur = cur.Advance()
		}

		c, ok := cur.Head()
		if !ok || c < '0' || c > '9' {
			return 0, in, ErrNoMatch
		}
		for ok && c >= '0' && c <= '9' {
			b.WriteByte(byte(c))
			cur = cur.Advance()
			c, ok = cur.Head()
		}

		if ok && c == '.' {
			b.WriteByte('.')
			cur = cur.Advance()
			c, ok = cur.Head()
			if !ok || c < '0' || c > '9' {
				return 0, in, NewError(cur, "expected digit after '.'")
			}
			for ok && c >= '0' && c <= '9' {
				b.WriteByte(byte(c))
				cur = cur.Advance()
				c, ok = cur.Head()
			}
		}

		if ok && (c == 'e' || c == 'E') {
			b.WriteByte(byte(c))
			cur = cur.Advance()
			c, ok = cur.Head()
			if ok && (c == '+' || c == '-') {
				b.WriteByte(byte(c))
				cur = cur.Advance()
				c, ok = cur.Head()
			}
			if !ok || c < '0' || c > '9' {
				return 0, in, NewError(cur, "expected digit in exponent")
			}
			for ok && c >= '0' && c <= '9' {
				b.WriteByte(byte(c))
				cur = cur.Advance()
				c, ok = cur.Head()
			}
		}

		f, err := strconv.ParseFloat(b.String(), 64)
		if err != nil {
			return 0, in, NewError(in, err.Error())
		}
		return f, cur, nil
	}
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
				return struct{}{}, in, NewError(in, fmt.Sprintf("unexpected %q", c))
			}
			return struct{}{}, in, NewError(in, "unexpected input")
		}
		return struct{}{}, in, nil
	}
}

// Label attaches a description to p, replacing its error message on failure.
func Label[T any](p Parser[T], label string) Parser[T] {
	return func(in Input) (T, Input, error) {
		val, next, err := p(in)
		if err != nil {
			return val, in, NewError(in, "expected "+label)
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
// Panics if p succeeds without consuming input, which would cause an infinite loop.
func ManyTill[T, E any](p Parser[T], end Parser[E]) Parser[[]T] {
	return func(in Input) ([]T, Input, error) {
		results := []T{}
		cur := in
		for {
			if _, next, err := end(cur); err == nil {
				return results, next, nil
			}
			val, next, err := p(cur)
			if err != nil {
				return nil, in, err
			}
			if next.Pos() == cur.Pos() {
				panic("parsec: ManyTill: parser succeeded without consuming input")
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
	return func(in Input) (int, Input, error) {
		c, ok := in.Head()
		if !ok || c < '0' || c > '9' {
			return 0, in, ErrNoMatch
		}
		n, cur := int(c-'0'), in.Advance()
		for {
			c, ok = cur.Head()
			if !ok || c < '0' || c > '9' {
				return n, cur, nil
			}
			n, cur = n*10+int(c-'0'), cur.Advance()
		}
	}
}
