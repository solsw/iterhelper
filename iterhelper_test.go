package iterhelper

import (
	"iter"
	"slices"
	"testing"

	"github.com/solsw/go2linq/v4"
)

func TestSeq2ToSeqK(t *testing.T) {
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
			got, err := Seq2ToSeqK(tt.args.seq2)
			if (err != nil) != tt.wantErr {
				t.Errorf("Seq2ToSeqK() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			equal, _ := go2linq.SequenceEqual(got, tt.want)
			if !equal {
				t.Errorf("Seq2ToSeqK() = %v, want %v", go2linq.StringDef(got), go2linq.StringDef(tt.want))
			}
		})
	}
}

func TestSeq2ToSeqV(t *testing.T) {
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
			got, err := Seq2ToSeqV(tt.args.seq2)
			if (err != nil) != tt.wantErr {
				t.Errorf("Seq2ToSeqV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			equal, _ := go2linq.SequenceEqual(got, tt.want)
			if !equal {
				t.Errorf("Seq2ToSeqV() = %v, want %v", go2linq.StringDef(got), go2linq.StringDef(tt.want))
			}
		})
	}
}
