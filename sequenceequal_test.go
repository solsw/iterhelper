package iterhelper

import (
	"iter"
	"testing"
)

func TestSequenceEqual_int(t *testing.T) {
	r0 := intSeq(0, 0)
	r1 := intSeq(0, 1)
	r2 := intSeq(0, 2)

	type args struct {
		first  iter.Seq[int]
		second iter.Seq[int]
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "EmptyEmpty",
			args: args{
				first:  Empty[int](),
				second: Empty[int](),
			},
			want: true,
		},
		{name: "EmptyFirst",
			args: args{
				first:  Empty[int](),
				second: VarSeq(2),
			},
			want: false,
		},
		{name: "EmptySecond",
			args: args{
				first:  VarSeq(1),
				second: Empty[int](),
			},
			want: false,
		},
		{name: "EqualSequences",
			args: args{
				first:  VarSeq(1),
				second: VarSeq(1),
			},
			want: true,
		},
		{name: "UnequalLengthsBothArrays",
			args: args{
				first:  VarSeq(1, 5, 3),
				second: VarSeq(1, 5, 3, 10),
			},
			want: false,
		},
		{name: "UnequalLengthsBothRangesFirstLonger",
			args: args{
				first:  intSeq(0, 11),
				second: intSeq(0, 10),
			},
			want: false,
		},
		{name: "UnequalLengthsBothRangesSecondLonger",
			args: args{
				first:  intSeq(0, 10),
				second: intSeq(0, 11),
			},
			want: false,
		},
		{name: "UnequalData",
			args: args{
				first:  VarSeq(1, 5, 3, 9),
				second: VarSeq(1, 5, 3, 10),
			},
			want: false,
		},
		{name: "EqualDataBothArrays",
			args: args{
				first:  VarSeq(1, 5, 3, 10),
				second: VarSeq(1, 5, 3, 10),
			},
			want: true,
		},
		{name: "EqualDataBothRanges",
			args: args{
				first:  intSeq(0, 10),
				second: intSeq(0, 10),
			},
			want: true,
		},
		{name: "OrderMatters",
			args: args{
				first:  VarSeq(1, 2),
				second: VarSeq(2, 1),
			},
			want: false,
		},
		{name: "Same0",
			args: args{
				first:  r0,
				second: r0,
			},
			want: true,
		},
		{name: "Same1",
			args: args{
				first:  r1,
				second: r1,
			},
			want: true,
		},
		{name: "Same2",
			args: args{
				first:  r2,
				second: r2,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := SequenceEqual(tt.args.first, tt.args.second)
			if got != tt.want {
				t.Errorf("SequenceEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSequenceEqual_string(t *testing.T) {
	type args struct {
		first  iter.Seq[string]
		second iter.Seq[string]
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "2",
			args: args{
				first:  VarSeq("one", "two", "three", "four"),
				second: VarSeq("one", "two", "three", "four"),
			},
			want: true,
		},
		{name: "4",
			args: args{
				first:  VarSeq("a", "b"),
				second: VarSeq("a"),
			},
			want: false,
		},
		{name: "5",
			args: args{
				first:  VarSeq("a"),
				second: VarSeq("a", "b"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := SequenceEqual(tt.args.first, tt.args.second)
			if got != tt.want {
				t.Errorf("SequenceEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSequenceEqualEq_string(t *testing.T) {
	type args struct {
		first  iter.Seq[string]
		second iter.Seq[string]
		equal  func(string, string) bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "1",
			args: args{
				first:  VarSeq("a", "b"),
				second: VarSeq("a", "B"),
				equal:  caseInsensitiveEqual,
			},
			want: true,
		},
		{name: "CustomEqualityComparer",
			args: args{
				first:  VarSeq("foo", "BAR", "baz"),
				second: VarSeq("FOO", "bar", "Baz"),
				equal:  caseInsensitiveEqual,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := SequenceEqualEq(tt.args.first, tt.args.second, tt.args.equal)
			if got != tt.want {
				t.Errorf("SequenceEqualEq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSequenceEqual2_int_string(t *testing.T) {
	type args struct {
		first  iter.Seq2[int, string]
		second iter.Seq2[int, string]
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "EmptyEmpty",
			args: args{
				first:  Empty2[int, string](),
				second: Empty2[int, string](),
			},
			want: true,
		},
		{name: "EmptyFirst",
			args: args{
				first:  Empty2[int, string](),
				second: sec2_int_string(1),
			},
			want: false,
		},
		{name: "EmptySecond",
			args: args{
				first:  sec2_int_string(1),
				second: Empty2[int, string](),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := SequenceEqual2(tt.args.first, tt.args.second)
			if got != tt.want {
				t.Errorf("SequenceEqual2() = %v, want %v", got, tt.want)
			}
		})
	}
}
