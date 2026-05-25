package parsec_test

import (
	"strings"
	"testing"

	"github.com/ajiyoshi-vg/goparsec/parsec"
)

func TestExpect_success(t *testing.T) {
	p := parsec.Expect(parsec.Natural(), "number")
	got, err := parsec.RunString(p, "42")
	if err != nil {
		t.Fatal(err)
	}
	if got != 42 {
		t.Errorf("got %d, want 42", got)
	}
}

func TestExpect_softFailure(t *testing.T) {
	// ErrNoMatch → replaced with label message at start position
	p := parsec.Expect(parsec.Natural(), "number")
	_, err := parsec.RunString(p, "abc")
	pe, ok := err.(*parsec.ParseError)
	if !ok {
		t.Fatalf("expected *ParseError, got %T: %v", err, err)
	}
	if pe.Message != "number" {
		t.Errorf("message = %q, want \"number\"", pe.Message)
	}
	if pe.Col != 1 {
		t.Errorf("col = %d, want 1", pe.Col)
	}
}

func TestExpect_hardFailure_messageReplaced(t *testing.T) {
	// Hard failure (input consumed): message is replaced, position preserved
	p := parsec.Expect(parsec.Literal("hello"), "greeting")
	_, err := parsec.RunString(p, "helXo")
	pe, ok := err.(*parsec.ParseError)
	if !ok {
		t.Fatalf("expected *ParseError, got %T: %v", err, err)
	}
	if pe.Message != "greeting" {
		t.Errorf("message = %q, want \"greeting\"", pe.Message)
	}
}

func TestExpect_hardFailure_positionPreserved(t *testing.T) {
	// Position must be where the hard failure occurred (col 4), not the start (col 1)
	p := parsec.Expect(parsec.Literal("hello"), "greeting")
	_, err := parsec.RunString(p, "helXo")
	pe, ok := err.(*parsec.ParseError)
	if !ok {
		t.Fatalf("expected *ParseError, got %T: %v", err, err)
	}
	if pe.Col != 4 {
		t.Errorf("col = %d, want 4 (position of 'X')", pe.Col)
	}
}

// TestExpect_vs_Label_nested shows nesting behaviour difference:
// Label: inner hard error passes through outer Label unchanged  → inner wins
// Expect: outer Expect replaces inner hard error message        → outer wins
func TestExpect_vs_Label_nested(t *testing.T) {
	inner := parsec.Literal("hello")

	// Label nesting: outer Label cannot replace hard error from inner Label
	withLabel := parsec.Label(parsec.Label(inner, "inner"), "outer")
	_, err := parsec.RunString(withLabel, "helXo")
	pe, ok := err.(*parsec.ParseError)
	if !ok {
		t.Fatalf("Label: expected *ParseError, got %T", err)
	}
	if strings.Contains(pe.Message, "outer") {
		t.Errorf("Label: outer label should not win, got %q", pe.Message)
	}

	// Expect nesting: outer Expect replaces even hard errors → outer wins
	withExpect := parsec.Expect(parsec.Expect(inner, "inner"), "outer")
	_, err = parsec.RunString(withExpect, "helXo")
	pe, ok = err.(*parsec.ParseError)
	if !ok {
		t.Fatalf("Expect: expected *ParseError, got %T", err)
	}
	if pe.Message != "outer" {
		t.Errorf("Expect: outer label should win, got %q", pe.Message)
	}
}
