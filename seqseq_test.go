package iterhelper

import (
	"errors"
	"fmt"
	"iter"
	"slices"
	"testing"

	"github.com/solsw/generichelper"
	"github.com/solsw/go2linq/v4"
)

func TestSeqSeq2_string_int_string(t *testing.T) {
	type args struct {
		seq      iter.Seq[string]
		selector func(string) (int, string)
	}
	tests := []struct {
		name        string
		args        args
		want        iter.Seq2[int, string]
		wantErr     bool
		expectedErr error
	}{
		{name: "NilSource",
			args: args{
				seq: nil,
				selector: func(s string) (int, string) {
					rr := []rune(s)
					slices.Reverse(rr)
					return len(s), string(rr)
				},
			},
			wantErr:     true,
			expectedErr: ErrNilSource,
		},
		{name: "NilSelector",
			args: args{
				seq:      VarSeq("one", "two", "three", "four"),
				selector: nil,
			},
			wantErr:     true,
			expectedErr: ErrNilSelector,
		},
		{name: "Empty",
			args: args{
				seq: Empty[string](),
				selector: func(s string) (int, string) {
					rr := []rune(s)
					slices.Reverse(rr)
					return len(s), string(rr)
				},
			},
			want: Empty2[int, string](),
		},
		{name: "Regular",
			args: args{
				seq: VarSeq("one", "two", "three", "four"),
				selector: func(s string) (int, string) {
					rr := []rune(s)
					slices.Reverse(rr)
					return len(s), string(rr)
				},
			},
			want: VarSeq2(
				generichelper.NewTuple2(3, "eno"),
				generichelper.NewTuple2(3, "owt"),
				generichelper.NewTuple2(5, "eerht"),
				generichelper.NewTuple2(4, "ruof"),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SeqSeq2(tt.args.seq, tt.args.selector)
			if (err != nil) != tt.wantErr {
				t.Errorf("SeqSeq2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("SeqSeq2() error = %v, expectedErr %v", errors.Unwrap(err), tt.expectedErr)
				}
				return
			}
			equal, _ := go2linq.SequenceEqual2(got, tt.want)
			if !equal {
				t.Errorf("SeqSeq2() = %v, want %v", go2linq.StringDef2(got), go2linq.StringDef2(tt.want))
			}
		})
	}
}

func TestSeq2Seq_int_string_string(t *testing.T) {
	type args struct {
		seq2     iter.Seq2[int, string]
		selector func(int, string) string
	}
	tests := []struct {
		name        string
		args        args
		want        iter.Seq[string]
		wantErr     bool
		expectedErr error
	}{
		{name: "NilSource",
			args: args{
				seq2:     nil,
				selector: func(i int, s string) string { return fmt.Sprintf("%d%s%[1]d", i, s) },
			},
			wantErr:     true,
			expectedErr: ErrNilSource,
		},
		{name: "NilSelector",
			args: args{
				seq2: VarSeq2(
					generichelper.NewTuple2(1, "one"),
					generichelper.NewTuple2(2, "two"),
					generichelper.NewTuple2(3, "three"),
					generichelper.NewTuple2(4, "four"),
				),
				selector: nil,
			},
			wantErr:     true,
			expectedErr: ErrNilSelector,
		},
		{name: "Empty",
			args: args{
				seq2:     Empty2[int, string](),
				selector: func(i int, s string) string { return fmt.Sprintf("%d%s%[1]d", i, s) },
			},
			want: Empty[string](),
		},
		{name: "Regular",
			args: args{
				seq2: VarSeq2(
					generichelper.NewTuple2(1, "one"),
					generichelper.NewTuple2(2, "two"),
					generichelper.NewTuple2(3, "three"),
					generichelper.NewTuple2(4, "four"),
				),
				selector: func(i int, s string) string { return fmt.Sprintf("%d%s%[1]d", i, s) },
			},
			want: VarSeq("1one1", "2two2", "3three3", "4four4"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Seq2Seq(tt.args.seq2, tt.args.selector)
			if (err != nil) != tt.wantErr {
				t.Errorf("Seq2Seq() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("Seq2Seq() error = %v, expectedErr %v", errors.Unwrap(err), tt.expectedErr)
				}
				return
			}
			equal, _ := SequenceEqual(got, tt.want)
			if !equal {
				t.Errorf("Seq2Seq() = %v, want %v", StringDef(got), StringDef(tt.want))
			}
		})
	}
}

func TestSeq2SeqK_int_string(t *testing.T) {
	type args struct {
		seq2 iter.Seq2[int, string]
	}
	tests := []struct {
		name    string
		args    args
		want    iter.Seq[int]
		wantErr bool
	}{
		{name: "error",
			args: args{
				seq2: nil,
			},
			wantErr: true,
		},
		{name: "empty",
			args: args{
				seq2: slices.All([]string{}),
			},
			want:    slices.Values([]int{}),
			wantErr: false,
		},
		{name: "normal",
			args: args{
				seq2: slices.All([]string{"one", "two"}),
			},
			want:    slices.Values([]int{0, 1}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Seq2SeqK(tt.args.seq2)
			if (err != nil) != tt.wantErr {
				t.Errorf("Seq2SeqK() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			equal, _ := SequenceEqual(got, tt.want)
			if !equal {
				t.Errorf("Seq2SeqK() = %v, want %v", StringDef(got), StringDef(tt.want))
			}
		})
	}
}

func TestSeq2SeqV_int_string(t *testing.T) {
	type args struct {
		seq2 iter.Seq2[int, string]
	}
	tests := []struct {
		name    string
		args    args
		want    iter.Seq[string]
		wantErr bool
	}{
		{name: "error",
			args: args{
				seq2: nil,
			},
			wantErr: true,
		},
		{name: "empty",
			args: args{
				seq2: slices.All([]string{}),
			},
			want:    slices.Values([]string{}),
			wantErr: false,
		},
		{name: "normal",
			args: args{
				seq2: slices.All([]string{"one", "two"}),
			},
			want:    slices.Values([]string{"one", "two"}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Seq2SeqV(tt.args.seq2)
			if (err != nil) != tt.wantErr {
				t.Errorf("Seq2SeqV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			equal, _ := SequenceEqual(got, tt.want)
			if !equal {
				t.Errorf("Seq2SeqV() = %v, want %v", StringDef(got), StringDef(tt.want))
			}
		})
	}
}
