package parsec_test

import (
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
