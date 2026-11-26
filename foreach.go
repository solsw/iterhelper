package iterhelper

import (
	"context"
	"iter"
	"runtime"

	"github.com/solsw/errorhelper"
	"golang.org/x/sync/errgroup"
)

// ForEach sequentially performs a specified 'action' on each value yielded by the [iterator].
// If 'ctx' is canceled or 'action' returns a non-nil error,
// the operation is stopped and corresponding error is returned.
//
// [iterator]: https://pkg.go.dev/iter#Seq
func ForEach[V any](ctx context.Context, seq iter.Seq[V], action func(V) error) error {
	if seq == nil {
		return errorhelper.CallerError(ErrNilSeq)
	}
	if action == nil {
		return errorhelper.CallerError(ErrNilAction)
	}
	if err := ctx.Err(); err != nil {
		return errorhelper.CallerError(err)
	}
	for v := range seq {
		if ctx.Err() != nil {
			break
		}
		if err := action(v); err != nil {
			return errorhelper.CallerError(err)
		}
		if ctx.Err() != nil {
			break
		}
	}
	return errorhelper.CallerError(ctx.Err())
}

// ForEachConcurrent concurrently performs a specified 'action' on each value yielded by the [iterator].
// If 'ctx' is canceled or 'action' returns a non-nil error,
// the operation is stopped and corresponding error is returned.
//
// [iterator]: https://pkg.go.dev/iter#Seq
func ForEachConcurrent[V any](ctx context.Context, seq iter.Seq[V], action func(V) error) error {
	if seq == nil {
		return errorhelper.CallerError(ErrNilSeq)
	}
	if action == nil {
		return errorhelper.CallerError(ErrNilAction)
	}
	err := runConcurrent(ctx, func(yield func(func() error) bool) {
		for v := range seq {
			if !yield(func() error { return action(v) }) {
				return
			}
		}
	})
	return errorhelper.CallerError(err)
}

// ForEach2 sequentially performs a specified 'action' on each pair of values yielded by the [iterator].
// If 'ctx' is canceled or 'action' returns a non-nil error,
// the operation is stopped and corresponding error is returned.
//
// [iterator]: https://pkg.go.dev/iter#Seq2
func ForEach2[K, V any](ctx context.Context, seq2 iter.Seq2[K, V], action func(K, V) error) error {
	if seq2 == nil {
		return errorhelper.CallerError(ErrNilSeq2)
	}
	if action == nil {
		return errorhelper.CallerError(ErrNilAction)
	}
	if err := ctx.Err(); err != nil {
		return errorhelper.CallerError(err)
	}
	for k, v := range seq2 {
		if ctx.Err() != nil {
			break
		}
		if err := action(k, v); err != nil {
			return errorhelper.CallerError(err)
		}
		if ctx.Err() != nil {
			break
		}
	}
	return errorhelper.CallerError(ctx.Err())
}

// ForEachConcurrent2 concurrently performs a specified 'action' on each pair of values yielded by the [iterator].
// If 'ctx' is canceled or 'action' returns a non-nil error,
// the operation is stopped and corresponding error is returned.
//
// [iterator]: https://pkg.go.dev/iter#Seq2
func ForEachConcurrent2[K, V any](ctx context.Context, seq2 iter.Seq2[K, V], action func(K, V) error) error {
	if seq2 == nil {
		return errorhelper.CallerError(ErrNilSeq2)
	}
	if action == nil {
		return errorhelper.CallerError(ErrNilAction)
	}
	err := runConcurrent(ctx, func(yield func(func() error) bool) {
		for k, v := range seq2 {
			if !yield(func() error { return action(k, v) }) {
				return
			}
		}
	})
	return errorhelper.CallerError(err)
}

// runConcurrent stops requesting actions on cancellation and waits for active calls.
func runConcurrent(ctx context.Context, actions iter.Seq[func() error]) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	group, workerCtx := errgroup.WithContext(ctx)
	slots := make(chan struct{}, 8*runtime.NumCPU())
	actions(func(action func() error) bool {
		if workerCtx.Err() != nil {
			return false
		}
		select {
		case <-workerCtx.Done():
			return false
		case slots <- struct{}{}:
		}
		group.Go(func() error {
			defer func() { <-slots }()
			if workerCtx.Err() != nil {
				return nil
			}
			return action()
		})
		return workerCtx.Err() == nil
	})
	if err := group.Wait(); err != nil {
		return err
	}
	return ctx.Err()
}
