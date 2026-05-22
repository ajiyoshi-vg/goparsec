package parsec

// Value runs p (discarding its result) and returns v.
// Useful for operator parsers: parse the operator token, return the combining function.
func Value[T, S any](p Parser[S], v T) Parser[T] {
	return Map(p, func(S) T { return v })
}
