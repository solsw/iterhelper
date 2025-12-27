package iterhelper

import (
	"fmt"
	"iter"
	"slices"
	"strings"

	"github.com/solsw/errorhelper"
)

// Format defines formatting parameters.
type Format struct {
	// Left and right rims surround each element of a sequence.
	LeftRim, RightRim string
	// Element separator separates elements of a sequence.
	ElementSeparator string
	// Left and right edges surround the whole string.
	LeftEdge, RightEdge string
	// Value separator separates values in pair.
	ValueSeparator string
}

// DefaultFormat represents default formatting parameters used by [StringDef] and [StringDef2].
// Assign desired values, if needed.
var DefaultFormat = Format{
	LeftRim:          "",
	RightRim:         "",
	ElementSeparator: " ",
	LeftEdge:         "[",
	RightEdge:        "]",
	ValueSeparator:   ":",
}

// StringFmt returns string representation of a [sequence] of values
// by calling [fmt.Sprint] on each value yielded by the [iterator].
// If 'seq' is nil, empty string is returned.
//
// [sequence]: https://pkg.go.dev/iter#Seq
// [iterator]: https://pkg.go.dev/iter#Seq
func StringFmt[V any](seq iter.Seq[V], format Format) string {
	if seq == nil {
		return ""
	}
	var b strings.Builder
	for v := range seq {
		if b.Len() > 0 {
			b.WriteString(format.ElementSeparator)
		}
		b.WriteString(format.LeftRim + fmt.Sprint(v) + format.RightRim)
	}
	return format.LeftEdge + b.String() + format.RightEdge
}

// StringDef returns string representation of a [sequence] of values
// by calling [fmt.Sprint] on each value yielded by the [iterator]
// and using default formatting parameters.
// If 'seq' is nil, empty string is returned.
//
// [sequence]: https://pkg.go.dev/iter#Seq
// [iterator]: https://pkg.go.dev/iter#Seq
func StringDef[V any](seq iter.Seq[V]) string {
	return StringFmt(seq, DefaultFormat)
}

// StringFmt2 returns string representation of a [sequence] of pairs of values
// by calling [fmt.Sprint] on each value yielded by the [iterator]:
// If 'seq2' is nil, empty string is returned.
//
// [sequence]: https://pkg.go.dev/iter#Seq2
// [iterator]: https://pkg.go.dev/iter#Seq2
func StringFmt2[K, V any](seq2 iter.Seq2[K, V], format Format) string {
	if seq2 == nil {
		return ""
	}
	var b strings.Builder
	for k, v := range seq2 {
		if b.Len() > 0 {
			b.WriteString(format.ElementSeparator)
		}
		b.WriteString(format.LeftRim + fmt.Sprint(k) + format.ValueSeparator + fmt.Sprint(v) + format.RightRim)
	}
	return format.LeftEdge + b.String() + format.RightEdge
}

// StringDef2 returns string representation of a [sequence] of pairs of values
// by calling [fmt.Sprint] on each value yielded by the [iterator]
// and using default formatting parameters.
// If 'seq2' is nil, empty string is returned.
//
// [sequence]: https://pkg.go.dev/iter#Seq2
// [iterator]: https://pkg.go.dev/iter#Seq2
func StringDef2[K, V any](seq2 iter.Seq2[K, V]) string {
	return StringFmt2(seq2, DefaultFormat)
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
