package iterhelper

import (
	"iter"
)

// Collect2 returns a slice of values collected from the [iterator].
// Each pair of values from the [iterator] results in two values in the slice.
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
