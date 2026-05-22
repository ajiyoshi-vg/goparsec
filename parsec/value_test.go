package parsec_test

import (
	"testing"

	"github.com/ajiyoshi-vg/goparsec/parsec"
)

func TestValue_success(t *testing.T) {
	p := parsec.Value(parsec.Char('+'), 42)
	got, err := parsec.RunString(p, "+")
	if err != nil {
		t.Fatal(err)
	}
	if got != 42 {
		t.Errorf("got %d, want 42", got)
	}
}

func TestValue_fail(t *testing.T) {
	p := parsec.Value(parsec.Char('+'), 42)
	_, err := parsec.RunString(p, "-")
	if err == nil {
		t.Error("expected error")
	}
}

func TestValue_consumesInput(t *testing.T) {
	p := parsec.Then(parsec.Value(parsec.Char('+'), 42), parsec.Natural())
	got, err := parsec.RunString(p, "+1")
	if err != nil {
		t.Fatal(err)
	}
	if got != 1 {
		t.Errorf("got %d, want 1", got)
	}
}

func TestValue_withFunc(t *testing.T) {
	add := func(a, b int) int { return a + b }
	p := parsec.Value(parsec.Char('+'), add)
	fn, err := parsec.RunString(p, "+")
	if err != nil {
		t.Fatal(err)
	}
	if fn(3, 4) != 7 {
		t.Errorf("fn(3,4) = %d, want 7", fn(3, 4))
	}
}
