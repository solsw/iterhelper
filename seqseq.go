package iterhelper

import (
	"iter"

	"github.com/solsw/errorhelper"
)

// SeqSeq2 converts [iter.Seq] to [iter.Seq2].
func SeqSeq2[V, K2, V2 any](seq iter.Seq[V], selector func(V) (K2, V2)) (iter.Seq2[K2, V2], error) {
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

// Seq2Seq converts [iter.Seq2] to [iter.Seq].
func Seq2Seq[K, V, V2 any](seq2 iter.Seq2[K, V], selector func(K, V) V2) (iter.Seq[V2], error) {
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

// Seq2SeqK converts [iter.Seq2] to iterator over keys of 'seq2'.
func Seq2SeqK[K, V any](seq2 iter.Seq2[K, V]) (iter.Seq[K], error) {
	if seq2 == nil {
		return nil, errorhelper.CallerError(ErrNilSource)
	}
	return Seq2Seq(seq2, func(k K, _ V) K { return k })
}

// Seq2SeqV converts [iter.Seq2] to iterator over values of 'seq2'.
func Seq2SeqV[K, V any](seq2 iter.Seq2[K, V]) (iter.Seq[V], error) {
	if seq2 == nil {
		return nil, errorhelper.CallerError(ErrNilSource)
	}
	return Seq2Seq(seq2, func(_ K, v V) V { return v })
}
