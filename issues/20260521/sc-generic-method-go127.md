# TODO: SC[S].Tok[U] — Go 1.27 の generic methods で空白コンシューマを改善

## 背景

現在、空白スキップを各トークンパーサに適用するために局所クロージャを使っている：

```go
w := parsec.Spaces()
tok := func(c rune) parsec.Parser[rune] { return parsec.Then(w, parsec.Char(c)) }
```

Go の型システムが generic function values をサポートしないため、型ごとにクロージャを定義する必要がある。

## Go 1.27 で可能になること

generic methods（Go 1.27 正式搭載予定、現在 tip で `GOEXPERIMENT=genericmethods`）を使えば：

```go
type SC[S any] struct{ skip Parser[S] }

func NewSC[S any](skip Parser[S]) SC[S] {
    return SC[S]{skip: skip}
}

func (sc SC[S]) Tok[U any](p Parser[U]) Parser[U] {
    return Then(sc.skip, p)
}
```

使用側：

```go
sc := parsec.NewSC(parsec.Spaces())
sc.Tok(parsec.Char('('))    // Parser[rune]
sc.Tok(parsec.Integer())    // Parser[int]
sc.Tok(parsec.GoString())   // Parser[string]
```

型ごとのクロージャが不要になる。

## 参照

- Go issue #77273: https://github.com/golang/go/issues/77273
- 実装状況コメント: https://github.com/golang/go/issues/77273#issuecomment-4360427977

## 対応タイミング

Go 1.27 リリース後に `SC[S]` 型を `parsec` パッケージに追加する。
