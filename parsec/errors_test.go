package parsec_test

import (
	"strings"
	"testing"

	"github.com/ajiyoshi-vg/goparsec/parsec"
)

func TestParseError_lineCol(t *testing.T) {
	// line 1, col 1: mismatch at start
	_, err := parsec.RunString(parsec.Char('x'), "abc")
	pe, ok := err.(*parsec.ParseError)
	if !ok {
		t.Fatalf("expected *ParseError, got %T", err)
	}
	if pe.Line() != 1 || pe.Col() != 1 {
		t.Errorf("expected line 1, col 1; got line %d, col %d", pe.Line(), pe.Col())
	}

	// after "hello\n", next position is line 2, col 1
	p := parsec.Then(parsec.Literal("hello\n"), parsec.Char('x'))
	_, err = parsec.RunString(p, "hello\nworld")
	pe, ok = err.(*parsec.ParseError)
	if !ok {
		t.Fatalf("expected *ParseError, got %T", err)
	}
	if pe.Line() != 2 || pe.Col() != 1 {
		t.Errorf("expected line 2, col 1; got line %d, col %d", pe.Line(), pe.Col())
	}

	// error message includes line:col
	if !strings.Contains(err.Error(), "line 2") {
		t.Errorf("error %q does not mention line 2", err.Error())
	}
}
