package iterhelper

import (
	"errors"
	"iter"
	"slices"

	"github.com/solsw/errorhelper"
	"github.com/solsw/generichelper"
)

var (
	ErrNilAction   = errors.New("nil action")
	ErrNilEqual    = errors.New("nil equal")
	ErrNilSelector = errors.New("nil selector")
	ErrNilSource   = errors.New("nil source")
)

// VarToSeq returns an [iterator] over the [variadic] parameters.
//
// [iterator]: https://pkg.go.dev/iter#hdr-Iterators
// [variadic]: https://go.dev/ref/spec#Function_types
func VarToSeq[E any](s ...E) iter.Seq[E] {
	return slices.Values(s)
}

// VarToSeq2 returns an [iterator] over the [variadic] tuples.
//
// [iterator]: https://pkg.go.dev/iter#hdr-Iterators
// [variadic]: https://go.dev/ref/spec#Function_types
func VarToSeq2[K, V any](tt ...generichelper.Tuple2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, t := range tt {
			if !yield(t.Item1, t.Item2) {
				return
			}
		}
	}
}

// SeqToSeq2 converts [iter.Seq] to [iter.Seq2].
func SeqToSeq2[V, K2, V2 any](seq iter.Seq[V], selector func(V) (K2, V2)) (iter.Seq2[K2, V2], error) {
	if seq == nil {
		return nil, errorhelper.CallerError(ErrNilSource)
	}
	if selector == nil {
		return nil, errorhelper.CallerError(ErrNilSelector)
	}
	return func(yield func(K2, V2) bool) {
			for v := range seq {
				if !yield(selector(v)) {
					return
				}
			}
		},
		nil
}

// Seq2ToSeq converts [iter.Seq2] to [iter.Seq].
func Seq2ToSeq[K, V, V2 any](seq2 iter.Seq2[K, V], selector func(K, V) V2) (iter.Seq[V2], error) {
	if seq2 == nil {
		return nil, errorhelper.CallerError(ErrNilSource)
	}
	if selector == nil {
		return nil, errorhelper.CallerError(ErrNilSelector)
	}
	return func(yield func(V2) bool) {
			for k, v := range seq2 {
				if !yield(selector(k, v)) {
					return
				}
			}
		},
		nil
}

// Seq2ToSeqK converts [iter.Seq2] to iterator over keys of 'seq2'.
func Seq2ToSeqK[K, V any](seq2 iter.Seq2[K, V]) (iter.Seq[K], error) {
	if seq2 == nil {
		return nil, errorhelper.CallerError(ErrNilSource)
	}
	return Seq2ToSeq(seq2, func(k K, _ V) K { return k })
}

// Seq2ToSeqV converts [iter.Seq2] to iterator over values of 'seq2'.
func Seq2ToSeqV[K, V any](seq2 iter.Seq2[K, V]) (iter.Seq[V], error) {
	if seq2 == nil {
		return nil, errorhelper.CallerError(ErrNilSource)
	}
	return Seq2ToSeq(seq2, func(_ K, v V) V { return v })
}
