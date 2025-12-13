package iterhelper

import (
	"iter"
	"testing"
)

func TestEqual_int(t *testing.T) {
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
				second: Var(2),
			},
			want: false,
		},
		{name: "EmptySecond",
			args: args{
				first:  Var(1),
				second: Empty[int](),
			},
			want: false,
		},
		{name: "EqualSequences",
			args: args{
				first:  Var(1),
				second: Var(1),
			},
			want: true,
		},
		{name: "UnequalLengthsBothArrays",
			args: args{
				first:  Var(1, 5, 3),
				second: Var(1, 5, 3, 10),
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
				first:  Var(1, 5, 3, 9),
				second: Var(1, 5, 3, 10),
			},
			want: false,
		},
		{name: "EqualDataBothArrays",
			args: args{
				first:  Var(1, 5, 3, 10),
				second: Var(1, 5, 3, 10),
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
				first:  Var(1, 2),
				second: Var(2, 1),
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
			got, _ := Equal(tt.args.first, tt.args.second)
			if got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEqual_string(t *testing.T) {
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
				first:  Var("one", "two", "three", "four"),
				second: Var("one", "two", "three", "four"),
			},
			want: true,
		},
		{name: "4",
			args: args{
				first:  Var("a", "b"),
				second: Var("a"),
			},
			want: false,
		},
		{name: "5",
			args: args{
				first:  Var("a"),
				second: Var("a", "b"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := Equal(tt.args.first, tt.args.second)
			if got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEqualEq_string(t *testing.T) {
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
				first:  Var("a", "b"),
				second: Var("a", "B"),
				equal:  caseInsensitiveEqual,
			},
			want: true,
		},
		{name: "CustomEqualityComparer",
			args: args{
				first:  Var("foo", "BAR", "baz"),
				second: Var("FOO", "bar", "Baz"),
				equal:  caseInsensitiveEqual,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := EqualEq(tt.args.first, tt.args.second, tt.args.equal)
			if got != tt.want {
				t.Errorf("EqualEq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEqual2_int_string(t *testing.T) {
	type args struct {
		first  iter.Seq2[int, string]
		second iter.Seq2[int, string]
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
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
		{name: "Equal",
			args: args{
				first:  sec2_int_string(2),
				second: sec2_int_string(2),
			},
			want: true,
		},
		{name: "NotEqual",
			args: args{
				first:  sec2_int_string(4),
				second: sec2_int_string(2),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := Equal2(tt.args.first, tt.args.second)
			if got != tt.want {
				t.Errorf("Equal2() = %v, want %v", got, tt.want)
			}
		})
	}
}
