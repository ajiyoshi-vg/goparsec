package parsec_test

// Integration test: S-expression arithmetic parser that builds an AST.
//
// Grammar:
//   sexp = integer | list
//   list = '(' op sexp+ ')'
//   op   = '+' | '-' | '*' | '/'

import (
	"testing"

	"github.com/ajiyoshi-vg/goparsec/input"
	"github.com/ajiyoshi-vg/goparsec/parsec"
)

// --- AST ---

type SexpNode interface{ sexpEval() int }

type SexpNum struct{ V int }

func (n SexpNum) sexpEval() int { return n.V }

type SexpForm struct {
	Op   rune
	Args []SexpNode
}

func (f SexpForm) sexpEval() int {
	switch f.Op {
	case '+':
		s := 0
		for _, a := range f.Args {
			s += a.sexpEval()
		}
		return s
	case '-':
		acc := f.Args[0].sexpEval()
		for _, a := range f.Args[1:] {
			acc -= a.sexpEval()
		}
		return acc
	case '*':
		p := 1
		for _, a := range f.Args {
			p *= a.sexpEval()
		}
		return p
	case '/':
		acc := f.Args[0].sexpEval()
		for _, a := range f.Args[1:] {
			acc /= a.sexpEval()
		}
		return acc
	}
	panic("unknown op: " + string(f.Op))
}

// --- Parser ---

func buildSexpASTParser() parsec.Parser[SexpNode] {
	var sexp parsec.Parser[SexpNode]
	w := parsec.Spaces()

	tok := func(c rune) parsec.Parser[rune] {
		return parsec.Then(w, parsec.Char(c))
	}

	sexpRef := parsec.Parser[SexpNode](func(in input.Input) (SexpNode, input.Input, error) {
		return sexp(in)
	})

	opP := parsec.Choice(tok('+'), tok('-'), tok('*'), tok('/'))

	listForm := parsec.Between(
		tok('('),
		parsec.Map2(opP, parsec.Many1(sexpRef), func(op rune, args []SexpNode) SexpNode {
			return SexpForm{Op: op, Args: args}
		}),
		tok(')'),
	)

	numP := parsec.Map(
		parsec.Then(w, parsec.Integer()),
		func(v int) SexpNode { return SexpNum{V: v} },
	)

	sexp = parsec.Choice(
		numP,
		parsec.Parser[SexpNode](func(in input.Input) (SexpNode, input.Input, error) {
			return listForm(in)
		}),
	)

	return sexp
}

// --- Tests ---

func TestSexp_num(t *testing.T) {
	p := buildSexpASTParser()

	got, err := parsec.RunStringFull(p, "42")
	if err != nil {
		t.Fatal(err)
	}
	n, ok := got.(SexpNum)
	if !ok {
		t.Fatalf("got %T, want SexpNum", got)
	}
	if n.V != 42 {
		t.Errorf("V = %d, want 42", n.V)
	}
}

func TestSexp_simpleForm(t *testing.T) {
	p := buildSexpASTParser()

	got, err := parsec.RunStringFull(p, "(+ 1 2)")
	if err != nil {
		t.Fatal(err)
	}
	f, ok := got.(SexpForm)
	if !ok {
		t.Fatalf("got %T, want SexpForm", got)
	}
	if f.Op != '+' {
		t.Errorf("Op = %q, want '+'", f.Op)
	}
	if len(f.Args) != 2 {
		t.Fatalf("len(Args) = %d, want 2", len(f.Args))
	}
	if f.Args[0].(SexpNum).V != 1 {
		t.Errorf("Args[0] = %v, want SexpNum{1}", f.Args[0])
	}
	if f.Args[1].(SexpNum).V != 2 {
		t.Errorf("Args[1] = %v, want SexpNum{2}", f.Args[1])
	}
}

func TestSexp_nested(t *testing.T) {
	p := buildSexpASTParser()

	got, err := parsec.RunStringFull(p, "(+ 1 (* 2 3) 4)")
	if err != nil {
		t.Fatal(err)
	}
	f, ok := got.(SexpForm)
	if !ok {
		t.Fatalf("got %T, want SexpForm", got)
	}
	if f.Op != '+' {
		t.Errorf("Op = %q, want '+'", f.Op)
	}
	if len(f.Args) != 3 {
		t.Fatalf("len(Args) = %d, want 3", len(f.Args))
	}
	inner, ok := f.Args[1].(SexpForm)
	if !ok {
		t.Fatalf("Args[1] = %T, want SexpForm", f.Args[1])
	}
	if inner.Op != '*' {
		t.Errorf("inner.Op = %q, want '*'", inner.Op)
	}
}

func TestSexp_eval(t *testing.T) {
	p := buildSexpASTParser()

	tests := []struct {
		input string
		want  int
	}{
		{"1", 1},
		{"-3", -3},
		{"(+ 1 2)", 3},
		{"(+ 1 2 3)", 6},
		{"(* 2 3)", 6},
		{"(- 10 3)", 7},
		{"(- 10 3 2)", 5},
		{"(/ 24 6)", 4},
		{"(/ 24 2 3)", 4},
		{"(+ 1 (* 2 3) 4)", 11},
		{"(* (+ 1 2) (+ 3 4))", 21},
		{"(+ (* 2 3) (* 4 5))", 26},
		{"( + 1 2 )", 3},
	}

	for _, tt := range tests {
		got, err := parsec.RunStringFull(p, tt.input)
		if err != nil {
			t.Errorf("Sexp(%q): %v", tt.input, err)
			continue
		}
		if v := got.sexpEval(); v != tt.want {
			t.Errorf("Sexp(%q).Eval() = %d, want %d", tt.input, v, tt.want)
		}
	}
}
