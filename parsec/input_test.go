package parsec_test

import (
	"strings"
	"testing"

	"github.com/ajiyoshi-vg/goparsec/parsec"
)

func TestInput_Head(t *testing.T) {
	in := parsec.NewInput("abc")

	c, ok := in.Head()
	if !ok || c != 'a' {
		t.Errorf("Head() = (%q, %v), want ('a', true)", c, ok)
	}
}

func TestInput_Head_EOF(t *testing.T) {
	in := parsec.NewInput("")

	_, ok := in.Head()
	if ok {
		t.Error("Head() on empty input should return false")
	}
}

func TestInput_Advance(t *testing.T) {
	in := parsec.NewInput("abc")
	next := in.Advance()

	c, ok := next.Head()
	if !ok || c != 'b' {
		t.Errorf("after Advance, Head() = (%q, %v), want ('b', true)", c, ok)
	}
}

func TestInput_IsEOF(t *testing.T) {
	in := parsec.NewInput("a")
	if in.IsEOF() {
		t.Error("IsEOF() on non-empty input should be false")
	}

	end := in.Advance()
	if !end.IsEOF() {
		t.Error("IsEOF() after consuming all input should be true")
	}
}

func TestParseError_lineCol(t *testing.T) {
	// line 1, col 1: mismatch at start
	_, err := parsec.Run(parsec.Char('x'), "abc")
	pe, ok := err.(*parsec.ParseError)
	if !ok {
		t.Fatalf("expected *ParseError, got %T", err)
	}
	if pe.Line != 1 || pe.Col != 1 {
		t.Errorf("expected line 1, col 1; got line %d, col %d", pe.Line, pe.Col)
	}

	// after "hello\n", next position is line 2, col 1
	p := parsec.Then(parsec.Literal("hello\n"), parsec.Char('x'))
	_, err = parsec.Run(p, "hello\nworld")
	pe, ok = err.(*parsec.ParseError)
	if !ok {
		t.Fatalf("expected *ParseError, got %T", err)
	}
	if pe.Line != 2 || pe.Col != 1 {
		t.Errorf("expected line 2, col 1; got line %d, col %d", pe.Line, pe.Col)
	}

	// error message includes line:col
	if !strings.Contains(err.Error(), "line 2") {
		t.Errorf("error %q does not mention line 2", err.Error())
	}
}
