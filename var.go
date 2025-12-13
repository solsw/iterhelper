package iterhelper

import (
	"iter"
	"slices"

	"github.com/solsw/generichelper"
)

// Var returns an [iterator] over the [variadic] parameters.
//
// [iterator]: https://pkg.go.dev/iter#Seq
// [variadic]: https://go.dev/ref/spec#Function_types
func Var[E any](s ...E) iter.Seq[E] {
	return slices.Values(s)
}

// Var2 returns an [iterator] over the [variadic] tuples.
//
// [iterator]: https://pkg.go.dev/iter#Seq2
// [variadic]: https://go.dev/ref/spec#Function_types
func Var2[K, V any](tt ...generichelper.Tuple2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, t := range tt {
			if !yield(t.Item1, t.Item2) {
				return
			}
		}
	}
}
