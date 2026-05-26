# Map2/Map3/Map4 をメソッドに昇格する

## 背景

Go は現在、ジェネリック型にジェネリックメソッドを定義できない。
そのため `Map2`/`Map3`/`Map4` は現在トップレベル関数として提供している。

```go
// 現在
Map2(opP, Many1(sexpRef), func(op rune, args []SexpNode) SexpNode { ... })

// Go 1.27 以降（想定）
opP.Map2(Many1(sexpRef), func(op rune, args []SexpNode) SexpNode { ... })
```

## 対応方針

Go が generic method をサポートしたタイミングで、`Parser[T]` に以下を追加する：

```go
func (p Parser[T]) Map[U any](f func(T) U) Parser[U]
func (p Parser[T]) Map2[S, U any](p2 Parser[S], f func(T, S) U) Parser[U]
func (p Parser[T]) Map3[S2, S3, U any](p2 Parser[S2], p3 Parser[S3], f func(T, S2, S3) U) Parser[U]
func (p Parser[T]) Map4[S2, S3, S4, U any](p2 Parser[S2], p3 Parser[S3], p4 Parser[S4], f func(T, S2, S3, S4) U) Parser[U]
func (p Parser[T]) Bind[U any](f func(T) Parser[U]) Parser[U]
```

トップレベル関数はそのまま残し、メソッドは内部で委譲する形にすれば後方互換を保てる。

## 備考

`Tuple`/`Pair` 型は `Map2` メソッドがあれば不要。
`Zip` スタイル（変換を後回し）も `Map2` で代替できる。
