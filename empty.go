package iterhelper

import (
	"iter"
)

// Empty returns an [iterator] over an empty [sequence] of values.
//
// [iterator]: https://pkg.go.dev/iter#hdr-Iterators
// [sequence]: https://pkg.go.dev/iter#Seq
func Empty[Result any]() iter.Seq[Result] {
	return func(func(Result) bool) {}
}

// Empty2 returns an [iterator] over an empty [sequence] of pairs of values.
//
// [iterator]: https://pkg.go.dev/iter#hdr-Iterators
// [sequence]: https://pkg.go.dev/iter#Seq2
func Empty2[K, V any]() iter.Seq2[K, V] {
	return func(func(K, V) bool) {}
}
