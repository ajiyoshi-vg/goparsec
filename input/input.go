package input

// Input is the interface that parser input streams must implement.
// All methods must be pure — implementations should be immutable value types.
type Input interface {
	// Head returns the next rune and true, or (0, false) at end of input.
	Head() (rune, bool)
	// Advance returns a new Input with the position moved one rune forward.
	Advance() Input
	// IsEOF reports whether the input is exhausted.
	IsEOF() bool
	// Pos returns the rune offset from the start (used for furthest-error comparison).
	Pos() int
	// Line returns the 1-based line number of the current position.
	Line() int
	// Col returns the 1-based column number of the current position.
	Col() int
}
