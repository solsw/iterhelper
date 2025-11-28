package iterhelper

import (
	"fmt"
	"iter"
	"slices"
	"strings"

	"github.com/solsw/errorhelper"
)

// StringFmt returns string representation of a sequence:
//   - if 'seq' is nil, empty string is returned;
//   - if 'T' implements [fmt.Stringer], it is used to convert each element to string;
//   - 'sep' separates elements;
//   - 'lrim' and 'rrim' surround each element;
//   - 'ledge' and 'redge' surround the whole string.
func StringFmt[T any](seq iter.Seq[T], sep, lrim, rrim, ledge, redge string) string {
	if seq == nil {
		return ""
	}
	var b strings.Builder
	for t := range seq {
		if b.Len() > 0 {
			b.WriteString(sep)
		}
		b.WriteString(lrim + fmt.Sprint(t) + rrim)
	}
	return ledge + b.String() + redge
}

// StringDef returns string representation of a sequence using default formatting.
// (See [StringFmt]: 'sep' is set to space, 'lrim' and 'rrim' are empty strings,
// 'ledge' and 'redge' are set to "[" and "]".)
func StringDef[T any](seq iter.Seq[T]) string {
	return StringFmt(seq, " ", "", "", "[", "]")
}

// StringFmt2 returns string representation of a sequence:
//   - if 'seq2' is nil, empty string is returned;
//   - if 'K' or 'V' implements [fmt.Stringer], it is used to convert each element to string;
//   - 'psep' separates pair of values;
//   - 'esep' separates elements;
//   - 'lrim' and 'rrim' surround each element;
//   - 'ledge' and 'redge' surround the whole string.
func StringFmt2[K, V any](seq2 iter.Seq2[K, V], psep, esep, lrim, rrim, ledge, redge string) string {
	if seq2 == nil {
		return ""
	}
	var b strings.Builder
	for k, v := range seq2 {
		if b.Len() > 0 {
			b.WriteString(esep)
		}
		b.WriteString(lrim + fmt.Sprint(k) + psep + fmt.Sprint(v) + rrim)
	}
	return ledge + b.String() + redge
}

// StringDef2 returns string representation of a sequence using default formatting.
// (See [StringFmt2]: 'psep' is set to colon, 'esep' is set to space,
// 'lrim' and 'rrim' are empty strings, 'ledge' and 'redge' are set to "[" and "]".)
func StringDef2[K, V any](seq2 iter.Seq2[K, V]) string {
	return StringFmt2(seq2, ":", " ", "", "", "[", "]")
}

// StringSeq converts a sequence to a sequence of strings.
func StringSeq[T any](seq iter.Seq[T]) (iter.Seq[string], error) {
	if seq == nil {
		return nil, errorhelper.CallerError(ErrNilSource)
	}
	return func(yield func(string) bool) {
			for t := range seq {
				if !yield(fmt.Sprint(t)) {
					return
				}
			}
		},
		nil
}

// StringSlice returns a sequence contents as a slice of strings.
func StringSlice[T any](seq iter.Seq[T]) ([]string, error) {
	seqString, err := StringSeq(seq)
	if err != nil {
		return nil, errorhelper.CallerError(err)
	}
	return slices.Collect(seqString), nil
}
