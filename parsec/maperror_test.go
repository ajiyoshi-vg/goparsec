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

// customErr is a user-defined error type that implements Positioned.
type customErr struct {
	pos int
	msg string
}

func (e *customErr) Error() string { return e.msg }
func (e *customErr) Offset() int   { return e.pos }

// TestPositioned_customErrorInChoice verifies that a user-defined error implementing
// Positioned participates correctly in Choice's furthest-error tracking.
func TestPositioned_customErrorInChoice(t *testing.T) {
	// p1 fails at position 0 with a custom error at offset 5
	p1 := parsec.Parser[int](func(in input.Input) (int, input.Input, error) {
		return 0, in, &customErr{pos: 5, msg: "custom at 5"}
	})
	// p2 fails at position 0 with a custom error at offset 2
	p2 := parsec.Parser[int](func(in input.Input) (int, input.Input, error) {
		return 0, in, &customErr{pos: 2, msg: "custom at 2"}
	})

	_, err := parsec.RunString(parsec.Choice(p1, p2), "anything")
	if err == nil {
		t.Fatal("expected error")
	}
	// p1 reached furthest (offset 5), so its error should win
	if err.Error() != "custom at 5" {
		t.Errorf("got %q, want \"custom at 5\"", err.Error())
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
