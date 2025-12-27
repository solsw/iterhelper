package iterhelper

import (
	"iter"
	"reflect"
	"testing"

	"github.com/solsw/errorhelper"
	"github.com/solsw/generichelper"
)

func TestCollect2(t *testing.T) {
	tests := []struct {
		name string
		seq2 iter.Seq2[int, string]
		want []any
	}{
		{
			name: "Empty",
			seq2: Empty2[int, string](),
			want: nil,
		},
		{
			name: "NonEmpty",
			seq2: errorhelper.Must(Var2[int, string](1, "one", 2, "two")),
			want: []any{1, "one", 2, "two"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Collect2(tt.seq2)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Collect2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollect2Tuple(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		seq2 iter.Seq2[int, string]
		want []generichelper.Tuple2[int, string]
	}{
		{
			name: "Empty",
			seq2: Empty2[int, string](),
			want: nil,
		},
		{
			name: "NonEmpty",
			seq2: errorhelper.Must(Var2[int, string](1, "one", 2, "two")),
			want: []generichelper.Tuple2[int, string]{{Item1: 1, Item2: "one"}, {Item1: 2, Item2: "two"}},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Collect2Tuple(tt.seq2)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Collect2Tuple() = %v, want %v", got, tt.want)
			}
		})
	}
}
