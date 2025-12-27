package iterhelper

import (
	"iter"
	"slices"

	"github.com/solsw/errorhelper"
	"github.com/solsw/generichelper"
)

// Var returns an [iterator] over the [variadic] parameters/values.
//
// [iterator]: https://pkg.go.dev/iter#Seq
// [variadic]: https://go.dev/ref/spec#Function_types
func Var[V any](vv ...V) iter.Seq[V] {
	return slices.Values(vv)
}

// Var2Tuple returns an [iterator] over the [variadic] tuples of values.
//
// [iterator]: https://pkg.go.dev/iter#Seq2
// [variadic]: https://go.dev/ref/spec#Function_types
func Var2Tuple[K, V any](tt ...generichelper.Tuple2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, t := range tt {
			if !yield(t.Item1, t.Item2) {
				return
			}
		}
	}
}

// Var2 returns an [iterator] over the [variadic] parameters/values.
// There must be an even number of parameters.
// Parameters must consist of pairs of values of type 'K' and 'V'.
//
// [iterator]: https://pkg.go.dev/iter#Seq2
// [variadic]: https://go.dev/ref/spec#Function_types
func Var2[K, V any](vv ...any) (iter.Seq2[K, V], error) {
	if (len(vv) % 2) != 0 {
		return nil, errorhelper.CallerError(ErrOddValues)
	}
	tt := make([]generichelper.Tuple2[K, V], len(vv)/2)
	for i := 0; i < len(vv); i += 2 {
		k, ok := vv[i].(K)
		if !ok {
			return nil, errorhelper.CallerError(ErrWrongType(vv[i], k), vv[i])
		}
		v, ok := vv[i+1].(V)
		if !ok {
			return nil, errorhelper.CallerError(ErrWrongType(vv[i+1], v), vv[i+1])
		}
		tt[i/2] = generichelper.Tuple2[K, V]{Item1: k, Item2: v}
	}
	return Var2Tuple(tt...), nil
}
