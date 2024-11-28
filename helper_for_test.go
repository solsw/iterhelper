package iterhelper

import (
	"fmt"
	"iter"
	"strings"
)

var caseInsensitiveEqual = func(x, y string) bool {
	return strings.ToLower(x) == strings.ToLower(y)
}

func intSeq(start, count int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := range count {
			if !yield(start + i) {
				return
			}
		}
	}
}

func sec2_int_string(n int) iter.Seq2[int, string] {
	return func(yield func(int, string) bool) {
		for i := range n {
			if !yield(i, fmt.Sprint(i)) {
				return
			}
		}
	}
}
