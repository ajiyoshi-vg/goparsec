package parsec

import "fmt"

// Satisfy parses a single rune satisfying pred.
func Satisfy(pred func(rune) bool) Parser[rune] {
	return func(in Input) (rune, Input, error) {
		c, ok := in.Head()
		if !ok {
			return 0, in, NewError(in, "unexpected end of input")
		}
		if !pred(c) {
			return 0, in, ErrNoMatch
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
		for i, c := range []rune(s) {
			r, ok := cur.Head()
			if !ok {
				return "", in, NewError(cur, fmt.Sprintf("expected %q, got EOF", s))
			}
			if r != c {
				if i == 0 {
					// Clean soft failure: no characters consumed yet.
					return "", in, ErrNoMatch
				}
				// Partial match: report the deeper position for better error messages.
				return "", in, NewError(cur, fmt.Sprintf("expected %q", s))
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
			return struct{}{}, in, NewError(in, fmt.Sprintf("expected EOF, got %q", c))
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

// AlphaNum parses an ASCII letter or digit.
func AlphaNum() Parser[rune] {
	return Satisfy(func(r rune) bool {
		return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
	})
}

// HexDigit parses a hexadecimal digit [0-9a-fA-F].
func HexDigit() Parser[rune] {
	return Satisfy(func(r rune) bool {
		return (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
	})
}
