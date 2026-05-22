package parsec

import (
	"strings"

	"github.com/ajiyoshi-vg/goparsec/input"
)

// GoString parses a Go double-quoted string literal with escape sequences.
//
// Supported escapes:
//   \\ \" \n \t \r \a \b \f \v
//   \xNN          (hex byte)
//   \uNNNN        (Unicode code point, 4 hex digits)
//   \UNNNNNNNN    (Unicode code point, 8 hex digits)
func GoString() Parser[string] {
	return Between(Char('"'), ManyChars(goStringChar()), Char('"'))
}

func goStringChar() Parser[rune] {
	return Choice(
		goEscapeSeq(),
		Satisfy(func(r rune) bool { return r != '"' && r != '\\' }),
	)
}

func goEscapeSeq() Parser[rune] {
	return Bind(Char('\\'), func(rune) Parser[rune] {
		return Choice(
			Map(Char('\\'), func(rune) rune { return '\\' }),
			Map(Char('"'), func(rune) rune { return '"' }),
			Map(Char('n'), func(rune) rune { return '\n' }),
			Map(Char('t'), func(rune) rune { return '\t' }),
			Map(Char('r'), func(rune) rune { return '\r' }),
			Map(Char('a'), func(rune) rune { return '\a' }),
			Map(Char('b'), func(rune) rune { return '\b' }),
			Map(Char('f'), func(rune) rune { return '\f' }),
			Map(Char('v'), func(rune) rune { return '\v' }),
			goHexEscape(),
			goUnicodeEscape('u', 4),
			goUnicodeEscape('U', 8),
		)
	})
}

func goHexEscape() Parser[rune] {
	return Bind(Char('x'), func(rune) Parser[rune] {
		return Map(Count(2, HexDigit()), hexDigitsToRune)
	})
}

func goUnicodeEscape(prefix rune, n int) Parser[rune] {
	return Bind(Char(prefix), func(rune) Parser[rune] {
		return Map(Count(n, HexDigit()), hexDigitsToRune)
	})
}

func hexDigitsToRune(ds []rune) rune {
	var r rune
	for _, d := range ds {
		r = r<<4 | hexDigitVal(d)
	}
	return r
}

func hexDigitVal(r rune) rune {
	switch {
	case r >= '0' && r <= '9':
		return r - '0'
	case r >= 'a' && r <= 'f':
		return r - 'a' + 10
	default:
		return r - 'A' + 10
	}
}

// JSONString parses a JSON double-quoted string literal (RFC 8259).
//
// Supported escapes:
//
//	\\ \" \/ \b \f \n \r \t
//	\uNNNN              (Unicode code point, 4 hex digits)
//	\uHHHH\uLLLL        (UTF-16 surrogate pair for code points > U+FFFF)
func JSONString() Parser[string] {
	return Between(Char('"'), ManyChars(jsonStringChar()), Char('"'))
}

func jsonStringChar() Parser[rune] {
	return Choice(
		jsonEscapeSeq(),
		// U+0020 and above, excluding " and \
		Satisfy(func(r rune) bool { return r >= 0x20 && r != '"' && r != '\\' }),
	)
}

func jsonEscapeSeq() Parser[rune] {
	return Bind(Char('\\'), func(rune) Parser[rune] {
		return Choice(
			Map(Char('\\'), func(rune) rune { return '\\' }),
			Map(Char('"'), func(rune) rune { return '"' }),
			Map(Char('/'), func(rune) rune { return '/' }),
			Map(Char('b'), func(rune) rune { return '\b' }),
			Map(Char('f'), func(rune) rune { return '\f' }),
			Map(Char('n'), func(rune) rune { return '\n' }),
			Map(Char('r'), func(rune) rune { return '\r' }),
			Map(Char('t'), func(rune) rune { return '\t' }),
			jsonUnicodeEscape(),
		)
	})
}

func jsonUnicodeEscape() Parser[rune] {
	hex4 := Bind(Char('u'), func(rune) Parser[rune] {
		return Map(Count(4, HexDigit()), hexDigitsToRune)
	})
	return Bind(hex4, func(r rune) Parser[rune] {
		if r >= 0xD800 && r <= 0xDBFF {
			// High surrogate: must be followed by \uDC00–\uDFFF
			return Bind(Char('\\'), func(rune) Parser[rune] {
				return Bind(hex4, func(lo rune) Parser[rune] {
					if lo < 0xDC00 || lo > 0xDFFF {
						return func(in input.Input) (rune, input.Input, error) {
							return 0, in, NewError(in, "invalid surrogate pair: low surrogate expected")
						}
					}
					return Return(0x10000 + (r-0xD800)*0x400 + (lo-0xDC00))
				})
			})
		}
		return Return(r)
	})
}

// ManyChars runs p zero or more times, accumulating runes directly into a string.
// It is the string-accumulating equivalent of Many — avoiding the []rune intermediate
// allocation that Map(Many(p), string) would produce.
func ManyChars(p Parser[rune]) Parser[string] {
	return func(in input.Input) (string, input.Input, error) {
		var b strings.Builder
		cur := in
		for {
			r, next, err := p(cur)
			if err != nil {
				return b.String(), cur, nil
			}
			if next.Pos() == cur.Pos() {
				panic("parsec: ManyChars: parser succeeded without consuming input")
			}
			b.WriteRune(r)
			cur = next
		}
	}
}

// Many1Chars runs p one or more times, accumulating runes into a string.
// Fails if p does not match at least once.
func Many1Chars(p Parser[rune]) Parser[string] {
	return func(in input.Input) (string, input.Input, error) {
		r, cur, err := p(in)
		if err != nil {
			return "", in, err
		}
		var b strings.Builder
		b.WriteRune(r)
		for {
			r, next, err := p(cur)
			if err != nil {
				return b.String(), cur, nil
			}
			if next.Pos() == cur.Pos() {
				panic("parsec: Many1Chars: parser succeeded without consuming input")
			}
			b.WriteRune(r)
			cur = next
		}
	}
}
