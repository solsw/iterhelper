package iterhelper

import (
	"iter"
	"reflect"

	"github.com/solsw/errorhelper"
	"github.com/solsw/generichelper"
)

// Equal determines whether two [iterator]s yield the equal sequences
// by comparing their elements using [generichelper.DeepEqual].
//
// [iterator]: https://pkg.go.dev/iter#Seq
func Equal[V any](first, second iter.Seq[V]) (bool, error) {
	if first == nil || second == nil {
		return false, errorhelper.CallerError(ErrNilSec)
	}
	r, err := EqualEq(first, second, generichelper.DeepEqual[V])
	if err != nil {
		return false, errorhelper.CallerError(err)
	}
	return r, nil
}

// EqualEq determines whether two [iterator]s yield the equal sequences
// by comparing their elements using a specified function.
//
// [iterator]: https://pkg.go.dev/iter#Seq
func EqualEq[V any](first, second iter.Seq[V], equal func(V, V) bool) (bool, error) {
	if first == nil || second == nil {
		return false, errorhelper.CallerError(ErrNilSec)
	}
	if equal == nil {
		return false, errorhelper.CallerError(ErrNilEqual)
	}
	next1, stop1 := iter.Pull(first)
	defer stop1()
	next2, stop2 := iter.Pull(second)
	defer stop2()
	for {
		v1, ok1 := next1()
		v2, ok2 := next2()
		if ok1 != ok2 {
			return false, nil
		}
		// here ok1 and ok2 are either both true or both false
		if !ok1 {
			break
		}
		if !equal(v1, v2) {
			return false, nil
		}
	}
	return true, nil
}

// Equal2 determines whether two [iterator]s yield the equal sequences
// by comparing their elements using [reflect.DeepEqual].
//
// [iterator]: https://pkg.go.dev/iter#Seq2
func Equal2[K, V any](first, second iter.Seq2[K, V]) (bool, error) {
	if first == nil || second == nil {
		return false, errorhelper.CallerError(ErrNilSec2)
	}
	r, err := EqualEq2(first, second,
		func(k1 K, v1 V, k2 K, v2 V) bool {
			return reflect.DeepEqual(k1, k2) && reflect.DeepEqual(v1, v2)
		})
	if err != nil {
		return false, errorhelper.CallerError(err)
	}
	return r, nil
}

// EqualEq2 determines whether two [iterator]s yield the equal sequences
// by comparing their elements using a specified function.
//
// [iterator]: https://pkg.go.dev/iter#Seq2
func EqualEq2[K, V any](first, second iter.Seq2[K, V],
	equal func(k1 K, v1 V, k2 K, v2 V) bool) (bool, error) {
	if first == nil || second == nil {
		return false, errorhelper.CallerError(ErrNilSec2)
	}
	if equal == nil {
		return false, errorhelper.CallerError(ErrNilEqual)
	}
	next1, stop1 := iter.Pull2(first)
	defer stop1()
	next2, stop2 := iter.Pull2(second)
	defer stop2()
	for {
		k1, v1, ok1 := next1()
		k2, v2, ok2 := next2()
		if ok1 != ok2 {
			return false, nil
		}
		// here ok1 and ok2 are either both true or both false
		if !ok1 {
			break
		}
		if !(equal(k1, v1, k2, v2)) {
			return false, nil
		}
	}
	return true, nil
}
