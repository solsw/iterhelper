package iterhelper

import (
	"fmt"
	"iter"
	"reflect"
	"testing"
)

func TestStringFmt_int(t *testing.T) {
	type args struct {
		seq   iter.Seq[int]
		sep   string
		lrim  string
		rrim  string
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
				seq:   VarToSeq(1, 2, 3, 4),
				sep:   "-",
				lrim:  "<",
				rrim:  ">",
				ledge: "[",
				redge: "]",
			},
			want: "[<1>-<2>-<3>-<4>]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringFmt(tt.args.seq, tt.args.sep, tt.args.lrim, tt.args.rrim, tt.args.ledge, tt.args.redge); got != tt.want {
				t.Errorf("StringFmt() = %v, want %v", got, tt.want)
			}
		})
	}
}

type intStringer int

func (i intStringer) String() string {
	return fmt.Sprintf("%d+%d", i, i*i)
}

func TestStringDef(t *testing.T) {
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
				seq: VarToSeq(intStringer(1), intStringer(2), intStringer(3)),
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
		{name: "1",
			args: args{
				seq: VarToSeq(any(intStringer(1)), any(2), any(intStringer(3))),
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
				seq: VarToSeq(1, 2, 3),
			},
			want: VarToSeq("1", "2", "3"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := StringSeq(tt.args.seq)
			equal, _ := SequenceEqual(got, tt.want)
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
				seq: VarToSeq(any(1), any(intStringer(2)), any(3)),
			},
			want: VarToSeq("1", "2+4", "3"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := StringSeq(tt.args.seq)
			equal, _ := SequenceEqual(got, tt.want)
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
				seq: VarToSeq(1, 2, 3),
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
				seq: VarToSeq(intStringer(1), 2, 3),
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
