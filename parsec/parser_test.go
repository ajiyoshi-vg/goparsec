package parsec_test

import (
	"strings"
	"testing"

	"github.com/ajiyoshi-vg/goparsec/input"
	"github.com/ajiyoshi-vg/goparsec/parsec"
)

func TestRun_withInput(t *testing.T) {
	in := input.NewString("123abc")
	got, err := parsec.Run(parsec.Natural(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 123 {
		t.Fatalf("got %d, want 123", got)
	}
}

func TestRun_withReaderAt(t *testing.T) {
	in := input.NewReaderAt(strings.NewReader("1,2,3"))
	p := parsec.SepBy(parsec.Natural(), parsec.Char(','))
	got, err := parsec.Run(p, in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []int{1, 2, 3}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestRunFull_withInput(t *testing.T) {
	in := input.NewString("abc")
	_, err := parsec.RunFull(parsec.Many1Chars(parsec.Letter()), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunFull_withInput_extraInput(t *testing.T) {
	in := input.NewString("ab1")
	_, err := parsec.RunFull(parsec.Many1Chars(parsec.Letter()), in)
	if err == nil {
		t.Fatal("expected error for unconsumed input, got nil")
	}
}

func TestRunString_basic(t *testing.T) {
	got, err := parsec.RunString(parsec.Natural(), "42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 42 {
		t.Fatalf("got %d, want 42", got)
	}
}
