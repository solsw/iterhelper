# iterhelper
[![Go Reference](https://pkg.go.dev/badge/github.com/solsw/iterhelper.svg)](https://pkg.go.dev/github.com/solsw/iterhelper)
[![GitHub](https://img.shields.io/badge/github--green?logo=github)](https://github.com/solsw/iterhelper)

`iterhelper` is a collection of helpers for Go's [`iter`](https://pkg.go.dev/iter) package.
It provides convenient functions to create, convert, compare, consume and format
[`iter.Seq`](https://pkg.go.dev/iter#Seq) and [`iter.Seq2`](https://pkg.go.dev/iter#Seq2) iterators.

## Installation

```sh
go get github.com/solsw/iterhelper
```

```go
import "github.com/solsw/iterhelper"
```

## Overview

Most functions come in two flavors: one for single-value sequences
([`iter.Seq[V]`](https://pkg.go.dev/iter#Seq)) and one, suffixed with `2`, for
key/value-pair sequences ([`iter.Seq2[K, V]`](https://pkg.go.dev/iter#Seq2)).

Functions that may fail on bad input (nil iterator, nil callback, wrong type)
return an `error` from the package's sentinel error set. Constructor-style
functions that cannot fail (`Empty`, `Var`, `ChanAll`, ...) return only an iterator.

## API

### Creating iterators

| Function | Description |
| --- | --- |
| `Empty[V]() iter.Seq[V]` | Iterator over an empty sequence. |
| `Empty2[K, V]() iter.Seq2[K, V]` | Iterator over an empty sequence of pairs. |
| `Var[V](vv ...V) iter.Seq[V]` | Iterator over variadic values. |
| `Var2[K, V](vv ...any) (iter.Seq2[K, V], error)` | Iterator over variadic values supplied as alternating `K`, `V` pairs (must be an even count of the proper types). |
| `Var2Tuple[K, V](tt ...Tuple2[K, V]) iter.Seq2[K, V]` | Iterator over variadic `generichelper.Tuple2` pairs. |
| `ChanAll[E](c <-chan E) iter.Seq[E]` | Iterator over the elements of a channel (nil channel yields an empty sequence). |
| `ChanAll2[E](c <-chan E) iter.Seq2[int, E]` | Iterator over index/element pairs of a channel. |

### Converting iterators

| Function | Description |
| --- | --- |
| `SeqSeq2[V, K2, V2](seq, selector) (iter.Seq2[K2, V2], error)` | Convert `iter.Seq` to `iter.Seq2` using a selector that maps each value to a pair. |
| `Seq2Seq[K, V, V2](seq2, selector) (iter.Seq[V2], error)` | Convert `iter.Seq2` to `iter.Seq` using a selector that maps each pair to a value. |
| `Seq2SeqK[K, V](seq2) (iter.Seq[K], error)` | Iterator over the keys of an `iter.Seq2`. |
| `Seq2SeqV[K, V](seq2) (iter.Seq[V], error)` | Iterator over the values of an `iter.Seq2`. |

### Collecting iterators

| Function | Description |
| --- | --- |
| `Collect2[K, V](seq2) []any` | Collect each pair into two consecutive slice elements. Returns `nil` for a nil iterator. |
| `Collect2Tuple[K, V](seq2) []generichelper.Tuple2[K, V]` | Collect each pair into a `Tuple2`. Returns `nil` for a nil iterator. |

### Comparing iterators

| Function | Description |
| --- | --- |
| `Equal[V](first, second) (bool, error)` | Compare two sequences element-by-element with [`reflect.DeepEqual`](https://pkg.go.dev/reflect#DeepEqual). |
| `EqualEq[V](first, second, equal) (bool, error)` | Compare two sequences with a custom equality function. |
| `Equal2[K, V](first, second) (bool, error)` | Compare two pair-sequences with `reflect.DeepEqual`. |
| `EqualEq2[K, V](first, second, equal) (bool, error)` | Compare two pair-sequences with a custom equality function. |

### Iterating

| Function | Description |
| --- | --- |
| `ForEach[V](ctx, seq, action) error` | Run `action` on each value sequentially; stops on context cancellation or an action error. |
| `ForEachConcurrent[V](ctx, seq, action) error` | Run `action` on each value concurrently (bounded by `8 * runtime.NumCPU()`). |
| `ForEach2[K, V](ctx, seq2, action) error` | Run `action` on each pair sequentially. |
| `ForEachConcurrent2[K, V](ctx, seq2, action) error` | Run `action` on each pair concurrently. |

### Formatting and string conversion

| Function | Description |
| --- | --- |
| `StringFmt[V](seq, format) string` | String representation of a sequence using the given `Format`. |
| `StringDef[V](seq) string` | String representation using `DefaultFormat`. |
| `StringFmt2[K, V](seq2, format) string` | String representation of a pair-sequence using the given `Format`. |
| `StringDef2[K, V](seq2) string` | String representation of a pair-sequence using `DefaultFormat`. |
| `StringSeq[V](seq) (iter.Seq[string], error)` | Convert an iterator to an iterator of strings via [`fmt.Sprint`](https://pkg.go.dev/fmt#Sprint). |
| `StringSlice[V](seq) ([]string, error)` | Collect a sequence into a `[]string` via `fmt.Sprint`. |

`Format` controls how a sequence is rendered:

```go
type Format struct {
    LeftRim, RightRim   string // surround each element
    ElementSeparator    string // between elements
    LeftEdge, RightEdge string // surround the whole string
    ValueSeparator      string // between key and value in a pair
}
```

`DefaultFormat` (used by `StringDef`/`StringDef2`) renders as `[a b c]` and
`[k1:v1 k2:v2]`. It is a package-level variable — assign to it to change the
defaults, but note that concurrent modification is not safe.

### Errors

The package exposes sentinel errors for use with [`errors.Is`](https://pkg.go.dev/errors#Is):

`ErrOddValues`, `ErrNilAction`, `ErrNilEqual`, `ErrNilSeq`, `ErrNilSeq2`,
`ErrNilSelector`, `ErrWrongType`.

## Examples

### Create and format

```go
seq := iterhelper.Var(1, 2, 3)
fmt.Println(iterhelper.StringDef(seq)) // [1 2 3]
```

### Build a pair-sequence and convert it

```go
seq2, err := iterhelper.Var2[string, int]("a", 1, "b", 2)
if err != nil {
    log.Fatal(err)
}
fmt.Println(iterhelper.StringDef2(seq2)) // [a:1 b:2]

keys, _ := iterhelper.Seq2SeqK(seq2)
fmt.Println(iterhelper.StringDef(keys)) // [a b]
```

### Iterate over a channel

```go
c := make(chan int)
go func() {
    defer close(c)
    for i := range 3 {
        c <- i
    }
}()
for v := range iterhelper.ChanAll(c) {
    fmt.Println(v)
}
```

### Run an action with context

```go
err := iterhelper.ForEach(ctx, iterhelper.Var(1, 2, 3), func(v int) error {
    fmt.Println(v)
    return nil
})
```

### Compare sequences

```go
eq, err := iterhelper.Equal(iterhelper.Var(1, 2, 3), iterhelper.Var(1, 2, 3))
// eq == true, err == nil
```

## Documentation

Full API documentation is available on
[pkg.go.dev](https://pkg.go.dev/github.com/solsw/iterhelper).
