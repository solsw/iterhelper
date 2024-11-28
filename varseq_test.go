package iterhelper

import (
	"iter"
	"testing"
)

func TestVarSeq_int_1(t *testing.T) {
	t.Run("", func(t *testing.T) {
		next, stop := iter.Pull(VarSeq(1))
		defer stop()
		_, _ = next()
		_, got := next()
		want := false
		if got != want {
			t.Errorf("VarSeq_1() = %v, want %v", got, want)
		}
	})
}

func TestVarSeq_int_2(t *testing.T) {
	t.Run("", func(t *testing.T) {
		next, stop := iter.Pull(VarSeq(1, 2))
		defer stop()
		_, _ = next()
		got, _ := next()
		want := 2
		if got != want {
			t.Errorf("VarSeq_2() = %v, want %v", got, want)
		}
	})
}
