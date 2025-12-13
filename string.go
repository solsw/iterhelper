package iterhelper

import (
	"fmt"
	"iter"
	"slices"
	"strings"

	"github.com/solsw/errorhelper"
)

// StringFmt returns string representation of a sequence of values yielded by the [iterator]
// by calling [fmt.Sprint] on each value yielded by the [iterator]:
//   - if 'seq' is nil, empty string is returned;
//   - 'lrim' and 'rrim' surround each value;
//   - 'sep' separates values;
//   - 'ledge' and 'redge' surround the whole string.
//
// [iterator]: https://pkg.go.dev/iter#Seq
func StringFmt[V any](seq iter.Seq[V], lrim, rrim, sep, ledge, redge string) string {
	if seq == nil {
		return ""
	}
	var b strings.Builder
	for v := range seq {
		if b.Len() > 0 {
			b.WriteString(sep)
		}
		b.WriteString(lrim + fmt.Sprint(v) + rrim)
	}
	return ledge + b.String() + redge
}

// StringDef returns string representation of a sequence of values
// yielded by the [iterator] using default formatting.
// (See [StringFmt]: 'lrim' and 'rrim' are empty strings,
// 'sep' is set to space, 'ledge' and 'redge' are set to "[" and "]".)
//
// [iterator]: https://pkg.go.dev/iter#Seq
func StringDef[V any](seq iter.Seq[V]) string {
	return StringFmt(seq, "", "", " ", "[", "]")
}

// StringFmt2 returns string representation of a sequence of pairs of values yielded by the [iterator]
// by calling [fmt.Sprint] on each value yielded by the [iterator]:
//   - if 'seq2' is nil, empty string is returned;
//   - 'vsep' separates values in pair;
//   - 'lrim' and 'rrim' surround each pair;
//   - 'psep' separates pairs;
//   - 'ledge' and 'redge' surround the whole string.
//
// [iterator]: https://pkg.go.dev/iter#Seq2
func StringFmt2[K, V any](seq2 iter.Seq2[K, V], vsep, lrim, rrim, psep, ledge, redge string) string {
	if seq2 == nil {
		return ""
	}
	var b strings.Builder
	for k, v := range seq2 {
		if b.Len() > 0 {
			b.WriteString(psep)
		}
		b.WriteString(lrim + fmt.Sprint(k) + vsep + fmt.Sprint(v) + rrim)
	}
	return ledge + b.String() + redge
}

// StringDef2 returns string representation of a sequence of pairs
// of values yielded by the [iterator] using default formatting.
// (See [StringFmt2]: 'vsep' is set to colon, 'lrim' and 'rrim' are empty strings,
// 'psep' is set to space, 'ledge' and 'redge' are set to "[" and "]".)
//
// [iterator]: https://pkg.go.dev/iter#Seq2
func StringDef2[K, V any](seq2 iter.Seq2[K, V]) string {
	return StringFmt2(seq2, ":", "", "", " ", "[", "]")
}

// StringSeq converts an [iterator] to an [iterator] over strings
// by calling [fmt.Sprint] on each value yielded by the [iterator].
//
// [iterator]: https://pkg.go.dev/iter#Seq
func StringSeq[V any](seq iter.Seq[V]) (iter.Seq[string], error) {
	if seq == nil {
		return nil, errorhelper.CallerError(ErrNilSec)
	}
	return func(yield func(string) bool) {
			for v := range seq {
				if !yield(fmt.Sprint(v)) {
					return
				}
			}
		},
		nil
}

// StringSlice returns a sequence of values yielded by the [iterator] as a slice of strings.
//
// [iterator]: https://pkg.go.dev/iter#Seq
func StringSlice[V any](seq iter.Seq[V]) ([]string, error) {
	seqString, err := StringSeq(seq)
	if err != nil {
		return nil, errorhelper.CallerError(err)
	}
	return slices.Collect(seqString), nil
}
