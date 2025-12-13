package iterhelper

import (
	"context"
	"iter"

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
		return errorhelper.CallerError(ErrNilSec)
	}
	if action == nil {
		return errorhelper.CallerError(ErrNilAction)
	}
	for v := range seq {
		select {
		case <-ctx.Done():
			return errorhelper.CallerError(ctx.Err())
		default:
			if err := action(v); err != nil {
				return errorhelper.CallerError(err)
			}
		}
	}
	return nil
}

// ForEachConcurrent concurrently performs a specified 'action' on each value yielded by the [iterator].
// If 'ctx' is canceled or 'action' returns a non-nil error,
// the operation is stopped and corresponding error is returned.
//
// [iterator]: https://pkg.go.dev/iter#Seq
func ForEachConcurrent[V any](ctx context.Context, seq iter.Seq[V], action func(V) error) error {
	if seq == nil {
		return errorhelper.CallerError(ErrNilSec)
	}
	if action == nil {
		return errorhelper.CallerError(ErrNilAction)
	}
	g := new(errgroup.Group)
	for v := range seq {
		g.Go(func() error {
			select {
			case <-ctx.Done():
				return errorhelper.CallerError(ctx.Err())
			default:
				if err := action(v); err != nil {
					return errorhelper.CallerError(err)
				}
			}
			return nil
		})
	}
	return errorhelper.CallerError(g.Wait())
}

// ForEach2 sequentially performs a specified 'action' on each pair of values yielded by the [iterator].
// If 'ctx' is canceled or 'action' returns a non-nil error,
// the operation is stopped and corresponding error is returned.
//
// [iterator]: https://pkg.go.dev/iter#Seq2
func ForEach2[K, V any](ctx context.Context, seq2 iter.Seq2[K, V], action func(K, V) error) error {
	if seq2 == nil {
		return errorhelper.CallerError(ErrNilSec)
	}
	if action == nil {
		return errorhelper.CallerError(ErrNilAction)
	}
	for k, v := range seq2 {
		select {
		case <-ctx.Done():
			return errorhelper.CallerError(ctx.Err())
		default:
			if err := action(k, v); err != nil {
				return errorhelper.CallerError(err)
			}
		}
	}
	return nil
}

// ForEachConcurrent2 concurrently performs a specified 'action' on each pair of values yielded by the [iterator].
// If 'ctx' is canceled or 'action' returns a non-nil error,
// the operation is stopped and corresponding error is returned.
//
// [iterator]: https://pkg.go.dev/iter#Seq2
func ForEachConcurrent2[K, V any](ctx context.Context, seq2 iter.Seq2[K, V], action func(K, V) error) error {
	if seq2 == nil {
		return errorhelper.CallerError(ErrNilSec)
	}
	if action == nil {
		return errorhelper.CallerError(ErrNilAction)
	}
	g := new(errgroup.Group)
	for k, v := range seq2 {
		g.Go(func() error {
			select {
			case <-ctx.Done():
				return errorhelper.CallerError(ctx.Err())
			default:
				if err := action(k, v); err != nil {
					return errorhelper.CallerError(err)
				}
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return errorhelper.CallerError(err)
	}
	return nil
}
