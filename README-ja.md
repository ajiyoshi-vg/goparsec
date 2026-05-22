# goparsec

Go のジェネリクスを活用したパーサコンビネータライブラリです。

```
import "github.com/ajiyoshi-vg/goparsec/parsec"
```

## インストール

```
go get github.com/ajiyoshi-vg/goparsec
```

Go 1.21 以降が必要です（`range N` 構文と `testing.B.Loop` を使用）。

## クイックスタート

```go
package main

import (
    "fmt"
    "github.com/ajiyoshi-vg/goparsec/parsec"
)

func main() {
    p := parsec.SepBy(parsec.Natural(), parsec.Char(','))
    got, err := parsec.Run(p, "1,2,3")
    fmt.Println(got, err) // [1 2 3] <nil>
}
```

## 基本概念

`Parser[T]` は次のシグネチャを持つ関数です：

```go
type Parser[T any] func(Input) (T, Input, error)
```

`Input` を受け取り、型 `T` の値と残りの `Input`、エラーを返します。
`Input` はイミュータブルな値型であるため、バックトラックは暗黙的です。
失敗したパーサは元の `Input` をそのまま返すだけで、状態を巻き戻す必要はありません。

### パーサの実行

```go
// Run は p で s を解析します。p が消費しなかった残りの入力は無視されます。
func Run[T any](p Parser[T], s string) (T, error)

// RunFull は p で s を解析し、入力が残っていた場合は失敗します。
func RunFull[T any](p Parser[T], s string) (T, error)
```

## プリミティブ

| 関数 | マッチする対象 |
|---|---|
| `Char(c)` | ルーン `c` と完全一致 |
| `AnyChar()` | 任意の 1 文字 |
| `Literal(s)` | 文字列 `s` と完全一致 |
| `Satisfy(pred)` | `pred(r)` が true となるルーン |
| `Digit()` | `0`–`9` |
| `Letter()` | ASCII アルファベット `a`–`z`、`A`–`Z` |
| `AlphaNum()` | ASCII アルファベットまたは数字 |
| `HexDigit()` | `0`–`9`、`a`–`f`、`A`–`F` |
| `Space()` | スペース、タブ、改行、キャリッジリターン |
| `EOF()` | 入力の終端（`struct{}{}` を返して成功） |

## コンビネータ

### 順序結合

```go
Then(pa, pb)            // pa を実行し、次に pb を実行して pb の結果を返す
Skip(pa, pb)            // pa を実行し、次に pb を実行して pa の結果を返す
Between(open, p, end)   // open、p、end を順に実行して p の結果を返す
Bind(p, f)              // p を実行し、その結果を f に渡して得たパーサを実行する
```

### 繰り返し

```go
Many(p)           // 0 回以上。常に成功
Many1(p)          // 1 回以上。1 回もマッチしなければ失敗
ManyChars(p)      // 0 回以上のルーンを文字列として収集
Many1Chars(p)     // 1 回以上のルーンを文字列として収集
Count(n, p)       // ちょうど n 回
SepBy(p, sep)     // sep で区切られた p を 0 回以上
ManyTill(p, end)  // end が成功するまで p を 0 回以上繰り返す。end を消費する
```

### 選択・省略可能

```go
Choice(p1, p2, ...)  // 最初に成功したものを返す。最も深い位置のエラーを報告
Option(def, p)       // p、または p が失敗した場合は def
```

### 変換

```go
Map(p, f)    // p の結果に f を適用する
Return(v)    // 入力を消費せず v で常に成功する
```

### エラー整形

```go
Label(p, "名前")       // p が失敗したときのメッセージを "expected 名前" に置き換える
NotFollowedBy(p)      // p が失敗した場合のみ成功し、入力を消費しない
```

### 演算子

```go
Chainl1(p, op)  // op で区切られた 1 回以上の p を左結合で畳み込む
Chainr1(p, op)  // op で区切られた 1 回以上の p を右結合で畳み込む
Spaces()        // 0 個以上の空白文字を文字列として返す
```

## 数値パーサ

```go
Natural() Parser[int]     // 1 桁以上の十進数 → 非負整数
Integer() Parser[int]     // 省略可能な '-' に続く数字 → 符号付き整数
Float()   Parser[float64] // [-] 整数部 [. 小数部] [e[+-] 指数部] → float64
```

## 文字列パーサ

```go
// Go のダブルクォート文字列リテラル: \n \t \\ \" \xNN \uNNNN \UNNNNNNNN
GoString() Parser[string]

// JSON 文字列 (RFC 8259): \n \t \\ \" \/ \b \f \r \uNNNN + サロゲートペア
JSONString() Parser[string]
```

## 入力ソース

デフォルトでは、パーサは `string` に対して動作します：

```go
parsec.Run(p, "入力文字列")
parsec.RunFull(p, "入力文字列")
```

`io.ReaderAt`（`*os.File`、`*strings.Reader`、`*bytes.Reader` など）から、全内容を `[]rune` バッファに展開せずにパースするには：

```go
// NewReaderInput は r を元にした Input を返します。
// 内容はオンデマンドで読み込まれます — rune バッファは確保されません。
// パース中、元の r は有効な状態を維持する必要があります。
func NewReaderInput(r io.ReaderAt) Input
```

使用例：

```go
f, _ := os.Open("data.txt")
defer f.Close()
in := parsec.NewReaderInput(f)
got, _, err := parsec.SepBy(parsec.Natural(), parsec.Char(','))(in)
```

`stringInput`（`Run`/`RunFull` が使用）は `[]rune` スライスを事前に一度だけ確保します。`readerInput` は文字ごとの I/O と引き換えに O(パーサスタック深さ) のメモリで動作します — 全内容をバッファリングしたくない大きなファイルに適しています。

## エラー

`*ParseError` は行・列・メッセージを保持します：

```go
type ParseError struct {
    Pos     int    // ルーンオフセット（最遠エラーの比較に使用）
    Line    int    // 1 始まり
    Col     int    // 1 始まり
    Message string
}
```

`Choice` は各代替案のうち、入力の最も深い位置に到達したものからエラーを報告します。

## カスタムパーサの書き方

カスタムパーサは `func(Input) (T, Input, error)` というシグネチャを持つ関数です。
`Choice` や `Label` と正しく連携するために、次の 2 つのエラーコンストラクタを使用してください：

```go
// ErrNoMatch はアロケーションなしで「ここではマッチしない」を表すセンチネルです。
// このエラーを受け取った Choice は次の代替案を試みます。
var ErrNoMatch error

// NewError は現在位置に紐づいた *ParseError を返します。
// 次の代替案に黙って進んでほしくない失敗に使います。
func NewError(in Input, msg string) error

// NewErrorf は NewError の fmt.Sprintf 版です。
func NewErrorf(in Input, format string, args ...any) error
```

例 — リテラル `42` にマッチするパーサ：

```go
fortyTwo := func(in parsec.Input) (int, parsec.Input, error) {
    c, ok := in.Head()
    if !ok || c != '4' {
        return 0, in, parsec.ErrNoMatch // ソフト失敗: Choice が次を試みる
    }
    cur := in.Advance()
    c, ok = cur.Head()
    if !ok || c != '2' {
        return 0, in, parsec.NewError(cur, "expected '2' after '4'") // ハード失敗
    }
    return 42, cur.Advance(), nil
}

p := parsec.Choice(parsec.Parser[int](fortyTwo), parsec.Natural())
got, _ := parsec.Run(p, "99") // → 99 (fortyTwo は ErrNoMatch を返し、Choice が Natural を試みる)
```

## 例：四則演算パーサ

`+`、`-`、`*`、`/`、括弧、負数リテラル、トークン間の空白に対応した再帰下降パーサです。

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

    factor := func(in parsec.Input) (int, parsec.Input, error) {
        paren := parsec.Between(
            tok('('),
            parsec.Parser[int](func(in parsec.Input) (int, parsec.Input, error) {
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

got, _ := parsec.RunFull(buildExpr(), "(1 + 2) * -3") // → -9
```

## ライセンス

MIT
