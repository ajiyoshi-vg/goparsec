package parsec

// Many parses zero or more occurrences of p. Always succeeds.
// p must consume input on success to avoid an infinite loop.
func Many[T any](p Parser[T]) Parser[[]T] {
	return func(in Input) ([]T, Input, error) {
		var results []T
		cur := in
		for {
			val, next, err := p(cur)
			if err != nil {
				return results, cur, nil
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
func Choice[T any](ps ...Parser[T]) Parser[T] {
	return func(in Input) (T, Input, error) {
		var zero T
		var lastErr error
		for _, p := range ps {
			val, next, err := p(in)
			if err == nil {
				return val, next, nil
			}
			lastErr = err
		}
		return zero, in, lastErr
	}
}

// Try runs p with backtracking. In this implementation all parsers already
// backtrack on failure, so Try is an identity provided for API compatibility.
func Try[T any](p Parser[T]) Parser[T] {
	return p
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
	return func(in Input) (T, Input, error) {
		var zero T
		_, next, err := open(in)
		if err != nil {
			return zero, in, err
		}
		val, next, err := p(next)
		if err != nil {
			return zero, in, err
		}
		_, next, err = close(next)
		if err != nil {
			return zero, in, err
		}
		return val, next, nil
	}
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

// SepBy1 parses one or more occurrences of p separated by sep.
func SepBy1[T, S any](p Parser[T], sep Parser[S]) Parser[[]T] {
	return func(in Input) ([]T, Input, error) {
		val, cur, err := p(in)
		if err != nil {
			return nil, in, err
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

// Lexeme runs p then skips trailing whitespace.
func Lexeme[T any](p Parser[T]) Parser[T] {
	return Skip(p, Spaces())
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
