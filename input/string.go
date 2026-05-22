package input

type stringInput struct {
	src  []rune
	pos  int
	line int
	col  int
}

// NewString returns an Input over the given string.
func NewString(s string) Input {
	return stringInput{src: []rune(s), line: 1, col: 1}
}

func (in stringInput) Head() (rune, bool) {
	if in.pos >= len(in.src) {
		return 0, false
	}
	return in.src[in.pos], true
}

func (in stringInput) Advance() Input {
	if in.pos >= len(in.src) {
		return in
	}
	next := stringInput{src: in.src, pos: in.pos + 1}
	if in.src[in.pos] == '\n' {
		next.line = in.line + 1
		next.col = 1
	} else {
		next.line = in.line
		next.col = in.col + 1
	}
	return next
}

func (in stringInput) IsEOF() bool { return in.pos >= len(in.src) }
func (in stringInput) Pos() int    { return in.pos }
func (in stringInput) Line() int   { return in.line }
func (in stringInput) Col() int    { return in.col }
