package iterhelper

import (
	"errors"
	"iter"

	"github.com/solsw/errorhelper"
)

var ErrNilIterator = errors.New("nil iterator")

// Seq2ToSeqK converts [iter.Seq2] to iterator over keys of 'seq2'.
func Seq2ToSeqK[K, V any](seq2 iter.Seq2[K, V]) (iter.Seq[K], error) {
	if seq2 == nil {
		return nil, errorhelper.CallerError(ErrNilIterator)
	}
	return func(yield func(K) bool) {
			for k, _ := range seq2 {
				if !yield(k) {
					return
				}
			}
		},
		nil
}

// Seq2ToSeqV converts [iter.Seq2] to iterator over values of 'seq2'.
func Seq2ToSeqV[K, V any](seq2 iter.Seq2[K, V]) (iter.Seq[V], error) {
	if seq2 == nil {
		return nil, errorhelper.CallerError(ErrNilIterator)
	}
	return func(yield func(V) bool) {
			for _, v := range seq2 {
				if !yield(v) {
					return
				}
			}
		},
		nil
}
