package iterhelper

import (
	"iter"

	"github.com/solsw/generichelper"
)

// Collect2 returns a slice of values collected from the [iterator].
// Each pair of values yielded by the [iterator] results in two values in the slice.
// If 'seq2' is nil, nil is returned.
//
// [iterator]: https://pkg.go.dev/iter#Seq2
func Collect2[K, V any](seq2 iter.Seq2[K, V]) []any {
	if seq2 == nil {
		return nil
	}
	var r []any
	for k, v := range seq2 {
		r = append(r, k, v)
	}
	return r
}

// Collect2Tuple returns a slice of tuples of values collected from the [iterator].
// If 'seq2' is nil, nil is returned.
//
// [iterator]: https://pkg.go.dev/iter#Seq2
func Collect2Tuple[K, V any](seq2 iter.Seq2[K, V]) []generichelper.Tuple2[K, V] {
	if seq2 == nil {
		return nil
	}
	var r []generichelper.Tuple2[K, V]
	for k, v := range seq2 {
		r = append(r, generichelper.Tuple2[K, V]{Item1: k, Item2: v})
	}
	return r
}
