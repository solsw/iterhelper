package iterhelper

import (
	"fmt"
	"iter"
	"reflect"
	"testing"

	"github.com/solsw/generichelper"
)

func TestStringFmt_int(t *testing.T) {
	type args struct {
		seq   iter.Seq[int]
		lrim  string
		rrim  string
		sep   string
		ledge string
		redge string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "1",
			args: args{
				seq:   Var(1, 2, 3, 4),
				lrim:  "<",
				rrim:  ">",
				sep:   "-",
				ledge: "[",
				redge: "]",
			},
			want: "[<1>-<2>-<3>-<4>]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StringFmt(tt.args.seq, tt.args.lrim, tt.args.rrim, tt.args.sep, tt.args.ledge, tt.args.redge)
			if got != tt.want {
				t.Errorf("StringFmt() = %v, want %v", got, tt.want)
			}
		})
	}
}

type intStringer int

func (i intStringer) String() string {
	return fmt.Sprintf("%d+%d", i, i*i)
}

func TestStringDef_Stringer(t *testing.T) {
	type args struct {
		seq iter.Seq[intStringer]
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "1",
			args: args{
				seq: Var(intStringer(1), intStringer(2), intStringer(3)),
			},
			want: "[1+1 2+4 3+9]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringDef(tt.args.seq); got != tt.want {
				t.Errorf("StringDef() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringDef_any(t *testing.T) {
	type args struct {
		seq iter.Seq[any]
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "nil seq",
			args: args{
				seq: nil,
			},
			want: "",
		},
		{name: "1",
			args: args{
				seq: Var(any(intStringer(1)), any(2), any(intStringer(3))),
			},
			want: "[1+1 2 3+9]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringDef(tt.args.seq); got != tt.want {
				t.Errorf("StringDef() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringSeq_int(t *testing.T) {
	type args struct {
		seq iter.Seq[int]
	}
	tests := []struct {
		name string
		args args
		want iter.Seq[string]
	}{
		{name: "1",
			args: args{
				seq: Var(1, 2, 3),
			},
			want: Var("1", "2", "3"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := StringSeq(tt.args.seq)
			equal, _ := Equal(got, tt.want)
			if !equal {
				t.Errorf("StringSeq() = %v, want %v", StringDef(got), StringDef(tt.want))
			}
		})
	}
}

func TestStringSeq_any(t *testing.T) {
	type args struct {
		seq iter.Seq[any]
	}
	tests := []struct {
		name string
		args args
		want iter.Seq[string]
	}{
		{name: "1",
			args: args{
				seq: Var(any(1), any(intStringer(2)), any(3)),
			},
			want: Var("1", "2+4", "3"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := StringSeq(tt.args.seq)
			equal, _ := Equal(got, tt.want)
			if !equal {
				t.Errorf("StringSeq() = %v, want %v", StringDef(got), StringDef(tt.want))
			}
		})
	}
}

func TestStringSlice_int(t *testing.T) {
	type args struct {
		seq iter.Seq[int]
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "1",
			args: args{
				seq: Var(1, 2, 3),
			},
			want: []string{"1", "2", "3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := StringSlice(tt.args.seq)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringSlice_intStringer(t *testing.T) {
	type args struct {
		seq iter.Seq[intStringer]
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "1",
			args: args{
				seq: Var(intStringer(1), 2, 3),
			},
			want: []string{"1+1", "2+4", "3+9"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := StringSlice(tt.args.seq)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringDef2_int_string(t *testing.T) {
	type args struct {
		seq2 iter.Seq2[int, string]
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "nil seq2",
			args: args{
				seq2: nil,
			},
			want: "",
		},
		{name: "1",
			args: args{
				seq2: Var2Tuple(
					generichelper.Tuple2[int, string]{Item1: 1, Item2: "one"},
					generichelper.Tuple2[int, string]{Item1: 2, Item2: "two"},
					generichelper.Tuple2[int, string]{Item1: 3, Item2: "three"},
				),
			},
			want: "[1:one 2:two 3:three]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringDef2(tt.args.seq2); got != tt.want {
				t.Errorf("StringDef2() = %v, want %v", got, tt.want)
			}
		})
	}
}
