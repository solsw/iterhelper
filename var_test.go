package iterhelper

import (
	"iter"
	"testing"

	"github.com/solsw/generichelper"
)

func TestVar_int_1(t *testing.T) {
	t.Run("", func(t *testing.T) {
		next, stop := iter.Pull(Var(1))
		defer stop()
		_, _ = next()
		_, got := next()
		want := false
		if got != want {
			t.Errorf("Var_1() = %v, want %v", got, want)
		}
	})
}

func TestVar_int_2(t *testing.T) {
	t.Run("", func(t *testing.T) {
		next, stop := iter.Pull(Var(1, 2))
		defer stop()
		_, _ = next()
		got, _ := next()
		want := 2
		if got != want {
			t.Errorf("Var_2() = %v, want %v", got, want)
		}
	})
}

func TestVar2(t *testing.T) {
	tests := []struct {
		name    string
		vv      []any
		want    iter.Seq2[int, string]
		wantErr bool
	}{
		{name: "odd number",
			vv:      []any{1},
			wantErr: true,
		},
		{name: "wrong key type",
			vv:      []any{"one", "two"},
			wantErr: true,
		},
		{name: "wrong value type",
			vv:      []any{1, 2},
			wantErr: true,
		},
		{name: "normal",
			vv: []any{1, "one", 2, "two"},
			want: Var2Tuple(
				generichelper.Tuple2[int, string]{Item1: 1, Item2: "one"},
				generichelper.Tuple2[int, string]{Item1: 2, Item2: "two"},
			),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := Var2[int, string](tt.vv...)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Var2() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Var2() succeeded unexpectedly")
			}
			equal, _ := Equal2(got, tt.want)
			if !equal {
				t.Errorf("Var2() = %v, want %v", got, tt.want)
			}
		})
	}
}
