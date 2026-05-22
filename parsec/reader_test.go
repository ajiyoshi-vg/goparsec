package parsec_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ajiyoshi-vg/goparsec/input"
	"github.com/ajiyoshi-vg/goparsec/parsec"
)

func TestReaderInput_withParser(t *testing.T) {
	in := input.NewReaderAt(strings.NewReader("1,2,3"))

	p := parsec.SepBy(parsec.Natural(), parsec.Char(','))
	got, _, err := p(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []int{1, 2, 3}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d] = %d, want %d", i, got[i], want[i])
		}
	}
}

func TestReaderInput_bytes(t *testing.T) {
	in := input.NewReaderAt(bytes.NewReader([]byte("hello")))

	p := parsec.Many1Chars(parsec.Letter())
	got, _, err := p(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "hello" {
		t.Fatalf("got %q, want %q", got, "hello")
	}
}
