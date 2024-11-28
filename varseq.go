package iterhelper

import (
	"iter"
	"slices"

	"github.com/solsw/generichelper"
)

// VarSeq returns an [iterator] over the [variadic] parameters.
//
// [iterator]: https://pkg.go.dev/iter#hdr-Iterators
// [variadic]: https://go.dev/ref/spec#Function_types
func VarSeq[E any](s ...E) iter.Seq[E] {
	return slices.Values(s)
}

// VarSeq2 returns an [iterator] over the [variadic] tuples.
//
// [iterator]: https://pkg.go.dev/iter#hdr-Iterators
// [variadic]: https://go.dev/ref/spec#Function_types
func VarSeq2[K, V any](tt ...generichelper.Tuple2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, t := range tt {
			if !yield(t.Item1, t.Item2) {
				return
			}
		}
	}
}
