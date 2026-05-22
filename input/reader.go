package input

import (
	"io"
	"unicode/utf8"
)

// readerInput is an Input backed by an io.ReaderAt.
// Each value is immutable; backtracking is free.
// Memory usage is O(parser stack depth) — no rune buffer is held.
//
// The underlying io.ReaderAt must remain valid for the duration of parsing;
// reads are issued on every Head() and Advance() call.
type readerInput struct {
	r       io.ReaderAt
	bytePos int64
	runePos int
	line    int
	col     int
}

// NewReaderAt returns an Input that reads from r on demand.
// r must implement io.ReaderAt (e.g. *strings.Reader, *bytes.Reader, *os.File).
// The content of r is never buffered; each Head() call issues a ReadAt.
func NewReaderAt(r io.ReaderAt) Input {
	return readerInput{r: r, bytePos: 0, runePos: 0, line: 1, col: 1}
}

func (in readerInput) Head() (rune, bool) {
	var buf [utf8.UTFMax]byte
	n, _ := in.r.ReadAt(buf[:], in.bytePos)
	if n == 0 {
		return 0, false
	}
	r, size := utf8.DecodeRune(buf[:n])
	if r == utf8.RuneError && size <= 1 {
		return 0, false
	}
	return r, true
}

func (in readerInput) Advance() Input {
	var buf [utf8.UTFMax]byte
	n, _ := in.r.ReadAt(buf[:], in.bytePos)
	if n == 0 {
		return in // at EOF; idempotent
	}
	r, size := utf8.DecodeRune(buf[:n])
	if r == utf8.RuneError && size <= 1 {
		return in
	}
	next := readerInput{
		r:       in.r,
		bytePos: in.bytePos + int64(size),
		runePos: in.runePos + 1,
	}
	if r == '\n' {
		next.line = in.line + 1
		next.col = 1
	} else {
		next.line = in.line
		next.col = in.col + 1
	}
	return next
}

func (in readerInput) IsEOF() bool {
	var buf [1]byte
	n, _ := in.r.ReadAt(buf[:], in.bytePos)
	return n == 0
}

func (in readerInput) Pos() int  { return in.runePos }
func (in readerInput) Line() int { return in.line }
func (in readerInput) Col() int  { return in.col }
