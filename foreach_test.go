package iterhelper

import (
	"context"
	"errors"
	"iter"
	"reflect"
	"sync/atomic"
	"testing"
	"testing/synctest"

	"github.com/solsw/errorhelper"
)

var ErrTestError = errors.New("test error")

func TestForEachStopsIteration(t *testing.T) {
	type runFunc func(context.Context, iter.Seq[int], func(int) error) error
	withPairs := func(run func(context.Context, iter.Seq2[int, int], func(int, int) error) error) runFunc {
		return func(ctx context.Context, seq iter.Seq[int], action func(int) error) error {
			pairs, err := SeqSeq2(seq, func(v int) (int, int) { return v, v })
			if err != nil {
				return err
			}
			return run(ctx, pairs, func(_, v int) error { return action(v) })
		}
	}
	for _, tc := range []struct {
		name string
		run  runFunc
	}{
		{"ForEach", ForEach[int]},
		{"ForEach2", withPairs(ForEach2[int, int])},
		{"ForEachConcurrent", ForEachConcurrent[int]},
		{"ForEachConcurrent2", withPairs(ForEachConcurrent2[int, int])},
	} {
		t.Run(tc.name+"/already canceled", func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			cancel()
			seq := func(func(int) bool) { t.Error("entered iterator after cancellation") }
			err := tc.run(ctx, seq, func(int) error { return nil })
			if !errors.Is(err, context.Canceled) {
				t.Fatalf("got %v, want context.Canceled", err)
			}
		})
		t.Run(tc.name+"/last action cancels", func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			err := tc.run(ctx, Var(1), func(int) error { cancel(); return nil })
			if !errors.Is(err, context.Canceled) {
				t.Fatalf("got %v, want context.Canceled", err)
			}
		})
		for _, cause := range []error{ErrTestError, context.Canceled} {
			t.Run(tc.name+"/"+cause.Error(), func(t *testing.T) {
				synctest.Test(t, func(t *testing.T) {
					ctx, cancel := context.WithCancel(t.Context())
					defer cancel()
					seq := func(yield func(int) bool) {
						if yield(1) {
							synctest.Wait() // Let the first action finish before offering more work.
							if yield(2) {
								t.Error("iterator was asked to continue after the action stopped iteration")
							}
						}
					}
					err := tc.run(ctx, seq, func(v int) error {
						if v != 1 {
							t.Errorf("unexpected action for %d", v)
						}
						if cause == context.Canceled {
							cancel()
							return nil
						}
						return cause
					})
					if !errors.Is(err, cause) {
						t.Fatalf("got %v, want %v", err, cause)
					}
				})
			})
		}
	}
}

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
			expectedErr: ErrNilSeq,
		},
		{name: "02",
			args: args{
				ctx: context.Background(),
				seq: Var(1, 2, 3),
			},
			wantErr:     true,
			expectedErr: ErrNilAction,
		},
		{name: "03",
			args: args{
				ctx: ctx1,
				seq: Var(1, 2, 3),
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
				seq: Var(1, 2, 3),
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
				seq: Var(1, 2, 3),
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
				seq:    Var(1, 2, 3),
				action: func(int) error { return nil },
			},
			wantErr:     true,
			expectedErr: context.Canceled,
		},
		{name: "02",
			args: args{
				ctx: context.Background(),
				seq: Var(1, 2, 3),
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

func TestForEach2_int_string(t *testing.T) {
	var acc1 int
	ctx1, cancel := context.WithCancel(context.Background())
	type args struct {
		ctx    context.Context
		seq2   iter.Seq2[int, string]
		action func(int, string) error
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
				action: func(k int, _ string) error {
					acc1 += k * k
					return nil
				},
			},
			wantErr:     true,
			expectedErr: ErrNilSeq2,
		},
		{name: "02",
			args: args{
				ctx:  context.Background(),
				seq2: sec2_int_string(3),
			},
			wantErr:     true,
			expectedErr: ErrNilAction,
		},
		{name: "03",
			args: args{
				ctx:  ctx1,
				seq2: sec2_int_string(4),
				action: func(k int, _ string) error {
					if k == 2 {
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
				ctx:  context.Background(),
				seq2: sec2_int_string(4),
				action: func(k int, _ string) error {
					if k == 2 {
						return errorhelper.CallerError(ErrTestError)
					}
					acc1 += k * k
					return nil
				},
			},
			wantErr:     true,
			expectedErr: ErrTestError,
		},
		{name: "1",
			args: args{
				ctx:  context.Background(),
				seq2: sec2_int_string(4),
				action: func(k int, _ string) error {
					acc1 += k * k
					return nil
				},
			},
			want: 14, // 0*0 + 1*1 + 2*2 + 3*3
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc1 = 0
			err := ForEach2(tt.args.ctx, tt.args.seq2, tt.args.action)
			if (err != nil) != tt.wantErr {
				t.Errorf("ForEach2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("ForEach2() error = %v, expectedErr %v", err, tt.expectedErr)
				}
				return
			}
			if !reflect.DeepEqual(acc1, tt.want) {
				t.Errorf("ForEach2() = %v, want %v", acc1, tt.want)
			}
		})
	}
}

func TestForEachConcurrent2_int_string(t *testing.T) {
	var acc1 int64
	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()
	type args struct {
		ctx    context.Context
		seq2   iter.Seq2[int, string]
		action func(int, string) error
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
				seq2:   sec2_int_string(3),
				action: func(int, string) error { return nil },
			},
			wantErr:     true,
			expectedErr: context.Canceled,
		},
		{name: "02",
			args: args{
				ctx:  context.Background(),
				seq2: sec2_int_string(3),
				action: func(k int, _ string) error {
					if k == 2 {
						return errorhelper.CallerError(ErrTestError)
					}
					atomic.AddInt64(&acc1, int64(k*k))
					return nil
				},
			},
			wantErr:     true,
			expectedErr: ErrTestError,
		},
		{name: "1",
			args: args{
				ctx:  context.Background(),
				seq2: sec2_int_string(1000),
				action: func(k int, _ string) error {
					atomic.AddInt64(&acc1, int64(k*k))
					return nil
				},
			},
			want: 332833500, // sum of k*k for k in 0..999
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc1 = 0
			err := ForEachConcurrent2(tt.args.ctx, tt.args.seq2, tt.args.action)
			if (err != nil) != tt.wantErr {
				t.Errorf("ForEachConcurrent2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("ForEachConcurrent2() error = %v, expectedErr %v", err, tt.expectedErr)
				}
				return
			}
			if !reflect.DeepEqual(acc1, tt.want) {
				t.Errorf("ForEachConcurrent2() = %v, want %v", acc1, tt.want)
			}
		})
	}
}
