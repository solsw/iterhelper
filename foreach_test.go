package iterhelper

import (
	"context"
	"errors"
	"iter"
	"reflect"
	"sync/atomic"
	"testing"

	"github.com/solsw/errorhelper"
)

var ErrTestError = errors.New("test error")

func TestForEach_int(t *testing.T) {
	var acc1 int
	ctx1, cancel := context.WithCancel(context.Background())
	type args struct {
		ctx    context.Context
		seq    iter.Seq[int]
		action func(int) error
	}
	tests := []struct {
		name        string
		args        args
		want        int
		wantErr     bool
		expectedErr error
	}{
		{name: "01",
			args: args{
				ctx: context.Background(),
				action: func(i int) error {
					acc1 += i * i
					return nil
				},
			},
			wantErr:     true,
			expectedErr: ErrNilSource,
		},
		{name: "02",
			args: args{
				ctx: context.Background(),
				seq: VarToSeq(1, 2, 3),
			},
			wantErr:     true,
			expectedErr: ErrNilAction,
		},
		{name: "03",
			args: args{
				ctx: ctx1,
				seq: VarToSeq(1, 2, 3),
				action: func(i int) error {
					if i == 2 {
						cancel()
					}
					return nil
				},
			},
			wantErr:     true,
			expectedErr: context.Canceled,
		},
		{name: "04",
			args: args{
				ctx: context.Background(),
				seq: VarToSeq(1, 2, 3),
				action: func(i int) error {
					if i == 2 {
						return errorhelper.CallerError(ErrTestError)
					}
					acc1 += i * i
					return nil
				},
			},
			wantErr:     true,
			expectedErr: ErrTestError,
		},
		{name: "1",
			args: args{
				ctx: context.Background(),
				seq: VarToSeq(1, 2, 3),
				action: func(i int) error {
					acc1 += i * i
					return nil
				},
			},
			want: 14,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc1 = 0
			err := ForEach(tt.args.ctx, tt.args.seq, tt.args.action)
			if (err != nil) != tt.wantErr {
				t.Errorf("ForEach() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("ForEach() error = %v, expectedErr %v", err, tt.expectedErr)
				}
				return
			}
			if !reflect.DeepEqual(acc1, tt.want) {
				t.Errorf("ForEach() = %v, want %v", acc1, tt.want)
			}
		})
	}
}

func TestForEachConcurrent_int(t *testing.T) {
	var acc1 int64
	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()
	type args struct {
		ctx    context.Context
		seq    iter.Seq[int]
		action func(int) error
	}
	tests := []struct {
		name        string
		args        args
		want        int64
		wantErr     bool
		expectedErr error
	}{
		{name: "01",
			args: args{
				ctx:    canceledCtx,
				seq:    VarToSeq(1, 2, 3),
				action: func(int) error { return nil },
			},
			wantErr:     true,
			expectedErr: context.Canceled,
		},
		{name: "02",
			args: args{
				ctx: context.Background(),
				seq: VarToSeq(1, 2, 3),
				action: func(i int) error {
					if i == 2 {
						return errorhelper.CallerError(ErrTestError)
					}
					atomic.AddInt64(&acc1, int64(i*i))
					return nil
				},
			},
			wantErr:     true,
			expectedErr: ErrTestError,
		},
		{name: "1",
			args: args{
				ctx: context.Background(),
				seq: intSeq(1, 1000),
				action: func(i int) error {
					// acc1 += int64(i * i) // <- demonstrates race error
					atomic.AddInt64(&acc1, int64(i*i))
					return nil
				},
			},
			want: 333833500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc1 = 0
			err := ForEachConcurrent(tt.args.ctx, tt.args.seq, tt.args.action)
			if (err != nil) != tt.wantErr {
				t.Errorf("ForEachConcurrent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("ForEachConcurrent() error = %v, expectedErr %v", err, tt.expectedErr)
				}
				return
			}
			if !reflect.DeepEqual(acc1, tt.want) {
				t.Errorf("ForEachConcurrent() = %v, want %v", acc1, tt.want)
			}
		})
	}
}
