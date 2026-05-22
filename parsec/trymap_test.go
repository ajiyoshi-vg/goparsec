package parsec_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/ajiyoshi-vg/goparsec/parsec"
)

func TestTryMap_success(t *testing.T) {
	// ManyChars(Digit()) → string → strconv.Atoi
	p := parsec.TryMap(parsec.ManyChars(parsec.Digit()), strconv.Atoi)
	got, err := parsec.RunString(p, "123")
	if err != nil {
		t.Fatal(err)
	}
	if got != 123 {
		t.Errorf("got %d, want 123", got)
	}
}

func TestTryMap_parserFails(t *testing.T) {
	p := parsec.TryMap(parsec.ManyChars(parsec.Digit()), strconv.Atoi)
	_, err := parsec.RunString(p, "abc")
	// ManyChars succeeds with "" then Atoi("") fails
	if err == nil {
		t.Error("expected error")
	}
}

func TestTryMap_fFails(t *testing.T) {
	mustPositive := func(n int) (int, error) {
		if n <= 0 {
			return 0, fmt.Errorf("must be positive, got %d", n)
		}
		return n, nil
	}
	p := parsec.TryMap(parsec.Integer(), mustPositive)

	got, err := parsec.RunString(p, "42")
	if err != nil {
		t.Fatal(err)
	}
	if got != 42 {
		t.Errorf("got %d, want 42", got)
	}

	_, err = parsec.RunString(p, "-1")
	if err == nil {
		t.Error("expected error for non-positive")
	}
}

func TestTryMap_errorIsParseError(t *testing.T) {
	mustPositive := func(n int) (int, error) {
		if n <= 0 {
			return 0, fmt.Errorf("must be positive")
		}
		return n, nil
	}
	p := parsec.TryMap(parsec.Integer(), mustPositive)

	_, err := parsec.RunString(p, "-1")
	if _, ok := err.(*parsec.ParseError); !ok {
		t.Errorf("expected *ParseError, got %T: %v", err, err)
	}
}
