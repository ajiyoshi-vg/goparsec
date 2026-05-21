package parsec

import "fmt"

// Satisfy parses a single rune satisfying pred.
func Satisfy(pred func(rune) bool) Parser[rune] {
	return func(in Input) (rune, Input, error) {
		c, ok := in.Head()
		if !ok {
			return 0, in, newError(in, "unexpected end of input")
		}
		if !pred(c) {
			return 0, in, newError(in, fmt.Sprintf("unexpected %q", c))
		}
		return c, in.Advance(), nil
	}
}

// Char parses the exact rune c.
func Char(c rune) Parser[rune] {
	return Satisfy(func(r rune) bool { return r == c })
}

// AnyChar parses any single rune.
func AnyChar() Parser[rune] {
	return Satisfy(func(rune) bool { return true })
}

// String parses the exact string s.
func String(s string) Parser[string] {
	return func(in Input) (string, Input, error) {
		cur := in
		for _, c := range []rune(s) {
			r, ok := cur.Head()
			if !ok {
				return "", in, newError(cur, fmt.Sprintf("expected %q, got EOF", s))
			}
			if r != c {
				return "", in, newError(cur, fmt.Sprintf("expected %q", s))
			}
			cur = cur.Advance()
		}
		return s, cur, nil
	}
}

// EOF succeeds only at end of input.
func EOF() Parser[struct{}] {
	return func(in Input) (struct{}, Input, error) {
		if !in.IsEOF() {
			c, _ := in.Head()
			return struct{}{}, in, newError(in, fmt.Sprintf("expected EOF, got %q", c))
		}
		return struct{}{}, in, nil
	}
}

// Digit parses a decimal digit.
func Digit() Parser[rune] {
	return Satisfy(func(r rune) bool { return r >= '0' && r <= '9' })
}

// Letter parses an ASCII letter.
func Letter() Parser[rune] {
	return Satisfy(func(r rune) bool {
		return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
	})
}

// Space parses a single whitespace character.
func Space() Parser[rune] {
	return Satisfy(func(r rune) bool {
		return r == ' ' || r == '\t' || r == '\n' || r == '\r'
	})
}
