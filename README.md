# goparsec

A parser combinator library for Go, built around generics.

```
import "github.com/ajiyoshi-vg/goparsec/parsec"
```

## Installation

```
go get github.com/ajiyoshi-vg/goparsec
```

Requires Go 1.21 or later (uses `range N` syntax and `testing.B.Loop`).

## Quick start

```go
package main

import (
    "fmt"
    "github.com/ajiyoshi-vg/goparsec/parsec"
)

func main() {
    p := parsec.SepBy(parsec.Natural(), parsec.Char(','))
    got, err := parsec.RunString(p, "1,2,3")
    fmt.Println(got, err) // [1 2 3] <nil>
}
```

## Core concept

A `Parser[T]` is a function:

```go
type Parser[T any] func(input.Input) (T, input.Input, error)
```

It consumes some of `Input`, returning a value of type `T` and the remaining
`Input`. Because `Input` is an immutable value type, backtracking is implicit:
a parser that fails simply returns the original `Input` unchanged; no state
needs to be rolled back.

### Running a parser

```go
// Run executes p on in; remaining input is ignored.
func Run[T any](p Parser[T], in input.Input) (T, error)

// RunFull executes p on in and fails if any input remains.
func RunFull[T any](p Parser[T], in input.Input) (T, error)

// RunString and RunStringFull are shorthands for the common case of a string input.
func RunString[T any](p Parser[T], s string) (T, error)
func RunStringFull[T any](p Parser[T], s string) (T, error)
```

## Primitives

| Function | Matches |
|---|---|
| `Char(c)` | the exact rune `c` |
| `AnyChar()` | any single rune |
| `Literal(s)` | the exact string `s` |
| `Satisfy(pred)` | a rune where `pred(r)` is true |
| `Digit()` | `0`–`9` |
| `Letter()` | ASCII letter `a`–`z`, `A`–`Z` |
| `AlphaNum()` | ASCII letter or digit |
| `HexDigit()` | `0`–`9`, `a`–`f`, `A`–`F` |
| `Space()` | space, tab, newline, or carriage return |
| `EOF()` | end of input (succeeds with `struct{}{}`) |

## Combinators

### Sequencing

```go
Then(pa, pb)            // run pa then pb, return pb's result
Skip(pa, pb)            // run pa then pb, return pa's result
Between(open, p, end)   // run open, p, end; return p's result
Bind(p, f)              // run p, pass result to f, run the returned parser
```

### Repetition

```go
Many(p)           // zero or more; always succeeds
Many1(p)          // one or more; fails if p never matches
ManyChars(p)      // zero or more runes, collected into a string
Many1Chars(p)     // one or more runes, collected into a string
Count(n, p)       // exactly n times
SepBy(p, sep)     // zero or more p, separated by sep
ManyTill(p, end)  // zero or more p until end succeeds; consumes end
```

### Choice and optionality

```go
Choice(p1, p2, ...)  // first match wins; reports the furthest failure
Option(def, p)       // p, or def if p fails
```

### Transformation

```go
Map(p, f)              // apply f to p's result
Map2(p1, p2, f)        // run p1 and p2 in sequence, apply f to both results
Map3(p1, p2, p3, f)    // run three parsers in sequence, apply f to all results
Map4(p1, p2, p3, p4, f)
Return(v)              // always succeeds with v, consuming nothing
Value(p, v)            // run p (discarding result), return v
```

### Error shaping

```go
Label(p, "name")      // replace p's failure message with "expected name"
NotFollowedBy(p)      // succeed only if p fails; consume nothing
```

### Operators

```go
Chainl1(p, op)  // one or more p separated by op, left-associative
Chainr1(p, op)  // one or more p separated by op, right-associative
Spaces()        // zero or more whitespace characters, returned as string
```

## Numeric parsers

```go
Natural() Parser[int]     // one or more decimal digits → non-negative int
Integer() Parser[int]     // optional '-' then digits → signed int
Float()   Parser[float64] // [-] digits [. digits] [e[+-] digits] → float64
```

## String parsers

```go
// Go double-quoted string literal: \n \t \\ \" \xNN \uNNNN \UNNNNNNNN
GoString() Parser[string]

// JSON string (RFC 8259): \n \t \\ \" \/ \b \f \r \uNNNN + surrogate pairs
JSONString() Parser[string]
```

## Input sources

Two input constructors are provided:

```go
// NewString wraps a string as an Input (pre-allocates a []rune slice).
func NewString(s string) input.Input

// NewReaderAt returns an Input backed by r.
// Content is read on demand — no rune buffer is allocated.
// The underlying r must remain valid for the duration of parsing.
func NewReaderAt(r io.ReaderAt) input.Input
```

For the common case of parsing a string, use the `RunString`/`RunStringFull` shorthands.
For any `Input` (including `input.NewReaderAt`), use `Run`/`RunFull` directly:

```go
// String shorthand
got, err := parsec.RunString(p, "1,2,3")

// io.ReaderAt source
f, _ := os.Open("data.txt")
defer f.Close()
got, err := parsec.Run(p, input.NewReaderAt(f))
```

`input.NewString` pre-allocates a single `[]rune` slice and never reads from disk. `input.NewReaderAt` trades per-character I/O for O(parser-stack-depth) memory — useful for large files where you do not want to buffer the entire content.

## Errors

`parsec.ParseError` carries the line, column, and a message:

```go
type ParseError struct {
    Pos     int    // rune offset (used for furthest-error comparison)
    Line    int    // 1-based
    Col     int    // 1-based
    Message string
}
```

`Choice` automatically picks the error from whichever alternative reached the
furthest position in the input.

## Writing custom parsers

A custom parser is any function with the signature `func(input.Input) (T, input.Input, error)`.
Use the error constructors from `parsec` to integrate correctly with `Choice` and `Label`:

```go
// ErrNoMatch is a zero-allocation sentinel for "did not match here".
// Choice will try the next alternative when it sees this error.
var ErrNoMatch error

// NewError returns a *ParseError at the current position.
// Use this for failures that should not be silently retried.
func NewError(in input.Input, msg string) error

// NewErrorf is like NewError but formats the message via fmt.Sprintf.
func NewErrorf(in input.Input, format string, args ...any) error
```

Example — a parser that matches the literal `42`:

```go
fortyTwo := func(in input.Input) (int, input.Input, error) {
    c, ok := in.Head()
    if !ok || c != '4' {
        return 0, in, parsec.ErrNoMatch // soft failure: Choice will try next
    }
    cur := in.Advance()
    c, ok = cur.Head()
    if !ok || c != '2' {
        return 0, in, parsec.NewError(cur, "expected '2' after '4'") // hard failure
    }
    return 42, cur.Advance(), nil
}

p := parsec.Choice(parsec.Parser[int](fortyTwo), parsec.Natural())
got, _ := parsec.RunString(p, "99") // → 99 (fortyTwo returns ErrNoMatch, Choice falls through)
```

## Example: expression parser

The following builds a recursive expression parser supporting `+`, `-`, `*`,
`/`, parentheses, negative literals, and whitespace between tokens.

```go
func buildExpr() parsec.Parser[int] {
    var expr parsec.Parser[int]
    w := parsec.Spaces()

    tok := func(c rune) parsec.Parser[rune] {
        return parsec.Then(w, parsec.Char(c))
    }
    op := func(c rune, fn func(int, int) int) parsec.Parser[func(int, int) int] {
        return parsec.Map(tok(c), func(rune) func(int, int) int { return fn })
    }

    factor := func(in input.Input) (int, input.Input, error) {
        paren := parsec.Between(
            tok('('),
            parsec.Parser[int](func(in input.Input) (int, input.Input, error) {
                return expr(in)
            }),
            tok(')'),
        )
        return parsec.Choice(parsec.Then(w, parsec.Integer()), paren)(in)
    }

    term := parsec.Chainl1(parsec.Parser[int](factor),
        parsec.Choice(op('*', func(a, b int) int { return a * b }),
                      op('/', func(a, b int) int { return a / b })))
    expr = parsec.Chainl1(term,
        parsec.Choice(op('+', func(a, b int) int { return a + b }),
                      op('-', func(a, b int) int { return a - b })))
    return expr
}

got, _ := parsec.RunStringFull(buildExpr(), "(1 + 2) * -3") // → -9
```

## License

MIT
