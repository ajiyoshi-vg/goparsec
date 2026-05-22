package parsec

// Map2 runs p1 and p2 in sequence, passing both results to f.
func Map2[A, B, C any](p1 Parser[A], p2 Parser[B], f func(A, B) C) Parser[C] {
	return Bind(p1, func(a A) Parser[C] {
		return Map(p2, func(b B) C { return f(a, b) })
	})
}

// Map3 runs p1, p2, and p3 in sequence, passing all results to f.
func Map3[A, B, C, D any](p1 Parser[A], p2 Parser[B], p3 Parser[C], f func(A, B, C) D) Parser[D] {
	return Bind(p1, func(a A) Parser[D] {
		return Map2(p2, p3, func(b B, c C) D { return f(a, b, c) })
	})
}

// Map4 runs p1, p2, p3, and p4 in sequence, passing all results to f.
func Map4[A, B, C, D, E any](p1 Parser[A], p2 Parser[B], p3 Parser[C], p4 Parser[D], f func(A, B, C, D) E) Parser[E] {
	return Bind(p1, func(a A) Parser[E] {
		return Map3(p2, p3, p4, func(b B, c C, d D) E { return f(a, b, c, d) })
	})
}
