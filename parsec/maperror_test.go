package parsec_test

import (
	"strings"
	"testing"

	"github.com/ajiyoshi-vg/goparsec/input"
	"github.com/ajiyoshi-vg/goparsec/parsec"
)

func TestMapError_success(t *testing.T) {
	p := parsec.MapError(parsec.Natural(), func(_ input.Input, err error) error {
		return err // never called on success
	})
	got, err := parsec.RunString(p, "42")
	if err != nil {
		t.Fatal(err)
	}
	if got != 42 {
		t.Errorf("got %d, want 42", got)
	}
}

func TestMapError_softFailure(t *testing.T) {
	p := parsec.MapError(parsec.Natural(), func(in input.Input, _ error) error {
		return parsec.NewError(in, "want number")
	})
	_, err := parsec.RunString(p, "abc")
	pe, ok := err.(*parsec.ParseError)
	if !ok {
		t.Fatalf("expected *ParseError, got %T: %v", err, err)
	}
	if pe.Message != "want number" {
		t.Errorf("message = %q, want \"want number\"", pe.Message)
	}
}

func TestMapError_hardFailure(t *testing.T) {
	p := parsec.MapError(parsec.Literal("hello"), func(in input.Input, err error) error {
		return parsec.NewError(in, "want greeting")
	})
	_, err := parsec.RunString(p, "helXo")
	pe, ok := err.(*parsec.ParseError)
	if !ok {
		t.Fatalf("expected *ParseError, got %T: %v", err, err)
	}
	if pe.Message != "want greeting" {
		t.Errorf("message = %q, want \"want greeting\"", pe.Message)
	}
}

func TestMapError_passThrough(t *testing.T) {
	p := parsec.MapError(parsec.Literal("hello"), func(_ input.Input, err error) error {
		return err // unchanged
	})
	_, err := parsec.RunString(p, "helXo")
	if err == nil {
		t.Fatal("expected error")
	}
	// original error message should survive
	if strings.Contains(err.Error(), "greeting") {
		t.Errorf("unexpected message: %v", err)
	}
}

// TestMapError_implementsExpect shows that Expect-like behaviour can be
// implemented via MapError, preserving the hard failure position.
func TestMapError_implementsExpect(t *testing.T) {
	expect := func(p parsec.Parser[int], label string) parsec.Parser[int] {
		return parsec.MapError(p, func(in input.Input, err error) error {
			if pe, ok := err.(*parsec.ParseError); ok {
				return &parsec.ParseError{Pos: pe.Pos, Line: pe.Line, Col: pe.Col, Message: label}
			}
			return parsec.NewError(in, label)
		})
	}

	p := expect(parsec.Natural(), "number")
	_, err := parsec.RunString(p, "abc")
	pe, ok := err.(*parsec.ParseError)
	if !ok {
		t.Fatalf("expected *ParseError, got %T", err)
	}
	if pe.Message != "number" {
		t.Errorf("message = %q, want \"number\"", pe.Message)
	}
}
