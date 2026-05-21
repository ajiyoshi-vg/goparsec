package parsec

// GoString parses a Go double-quoted string literal with escape sequences.
//
// Supported escapes:
//   \\ \" \n \t \r \a \b \f \v
//   \xNN          (hex byte)
//   \uNNNN        (Unicode code point, 4 hex digits)
//   \UNNNNNNNN    (Unicode code point, 8 hex digits)
func GoString() Parser[string] {
	return Between(
		Char('"'),
		Char('"'),
		Map(Many(goStringChar()), func(rs []rune) string { return string(rs) }),
	)
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
