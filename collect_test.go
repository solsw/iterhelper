package iterhelper

import (
	"iter"
	"reflect"
	"testing"

	"github.com/solsw/errorhelper"
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
