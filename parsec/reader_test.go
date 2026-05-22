package parsec_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ajiyoshi-vg/goparsec/parsec"
)

func TestReaderInput_basic(t *testing.T) {
	r := strings.NewReader("abc")
	in := parsec.NewReaderAtInput(r)

	c, ok := in.Head()
	if !ok || c != 'a' {
		t.Fatalf("Head() = (%q, %v), want ('a', true)", c, ok)
	}
	if in.Pos() != 0 {
		t.Fatalf("Pos() = %d, want 0", in.Pos())
	}
	if in.Line() != 1 {
		t.Fatalf("Line() = %d, want 1", in.Line())
	}
	if in.Col() != 1 {
		t.Fatalf("Col() = %d, want 1", in.Col())
	}

	in2 := in.Advance()
	c, ok = in2.Head()
	if !ok || c != 'b' {
		t.Fatalf("after Advance: Head() = (%q, %v), want ('b', true)", c, ok)
	}
	if in2.Pos() != 1 {
		t.Fatalf("after Advance: Pos() = %d, want 1", in2.Pos())
	}
	if in2.Col() != 2 {
		t.Fatalf("after Advance: Col() = %d, want 2", in2.Col())
	}
}

func TestReaderInput_EOF(t *testing.T) {
	r := strings.NewReader("")
	in := parsec.NewReaderAtInput(r)

	if !in.IsEOF() {
		t.Fatal("IsEOF() = false, want true on empty input")
	}
	_, ok := in.Head()
	if ok {
		t.Fatal("Head() on empty input returned ok=true")
	}
}

func TestReaderInput_advancePastEnd(t *testing.T) {
	r := strings.NewReader("x")
	in := parsec.NewReaderAtInput(r)

	in2 := in.Advance()
	if !in2.IsEOF() {
		t.Fatal("expected EOF after advancing past last rune")
	}
	// Advance at EOF must be idempotent
	in3 := in2.Advance()
	if in3.Pos() != in2.Pos() {
		t.Fatalf("Advance at EOF changed Pos: %d -> %d", in2.Pos(), in3.Pos())
	}
}

func TestReaderInput_lineCol(t *testing.T) {
	// "a\nbc"  positions:
	//   'a'  line=1 col=1
	//   '\n' line=1 col=2  → after Advance: line=2 col=1
	//   'b'  line=2 col=1
	//   'c'  line=2 col=2
	r := strings.NewReader("a\nbc")
	in := parsec.NewReaderAtInput(r)

	in = in.Advance() // past 'a'
	in = in.Advance() // past '\n'

	c, ok := in.Head()
	if !ok || c != 'b' {
		t.Fatalf("Head() = (%q, %v), want ('b', true)", c, ok)
	}
	if in.Line() != 2 {
		t.Fatalf("Line() = %d, want 2", in.Line())
	}
	if in.Col() != 1 {
		t.Fatalf("Col() = %d, want 1", in.Col())
	}
}

func TestReaderInput_multibyte(t *testing.T) {
	// "日本語" — each rune is 3 bytes in UTF-8
	r := strings.NewReader("日本語")
	in := parsec.NewReaderAtInput(r)

	c, ok := in.Head()
	if !ok || c != '日' {
		t.Fatalf("Head() = (%q, %v), want ('日', true)", c, ok)
	}
	if in.Pos() != 0 {
		t.Fatalf("Pos() = %d, want 0", in.Pos())
	}

	in = in.Advance()
	c, ok = in.Head()
	if !ok || c != '本' {
		t.Fatalf("after Advance: Head() = (%q, %v), want ('本', true)", c, ok)
	}
	if in.Pos() != 1 {
		t.Fatalf("after Advance: Pos() = %d, want 1 (rune offset)", in.Pos())
	}

	in = in.Advance()
	c, ok = in.Head()
	if !ok || c != '語' {
		t.Fatalf("Head() = (%q, %v), want ('語', true)", c, ok)
	}
	if in.Pos() != 2 {
		t.Fatalf("Pos() = %d, want 2", in.Pos())
	}

	in = in.Advance()
	if !in.IsEOF() {
		t.Fatal("expected EOF after last rune")
	}
}

func TestReaderInput_backtrackImplicit(t *testing.T) {
	// Verify that the original Input is unchanged after Advance (immutable)
	r := strings.NewReader("abc")
	in := parsec.NewReaderAtInput(r)

	in2 := in.Advance()
	_ = in2.Advance()

	// in should still see 'a'
	c, ok := in.Head()
	if !ok || c != 'a' {
		t.Fatalf("original Input mutated: Head() = (%q, %v), want ('a', true)", c, ok)
	}
}

func TestReaderInput_withParser(t *testing.T) {
	r := strings.NewReader("1,2,3")
	in := parsec.NewReaderAtInput(r)

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
	// bytes.Reader also implements io.ReaderAt
	r := bytes.NewReader([]byte("hello"))
	in := parsec.NewReaderAtInput(r)

	p := parsec.Many1Chars(parsec.Letter())
	got, _, err := p(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "hello" {
		t.Fatalf("got %q, want %q", got, "hello")
	}
}
