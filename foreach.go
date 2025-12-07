package iterhelper

import (
	"context"
	"iter"

	"github.com/solsw/errorhelper"
	"golang.org/x/sync/errgroup"
)

// ForEach sequentially performs a specified 'action' on each element of the sequence.
// If 'ctx' is canceled or 'action' returns non-nil error,
// operation is stopped and corresponding error is returned.
func ForEach[T any](ctx context.Context, seq iter.Seq[T], action func(T) error) error {
	if seq == nil {
		return errorhelper.CallerError(ErrNilSec)
	}
	if action == nil {
		return errorhelper.CallerError(ErrNilAction)
	}
	for t := range seq {
		select {
		case <-ctx.Done():
			return errorhelper.CallerError(ctx.Err())
		default:
			if err := action(t); err != nil {
				return errorhelper.CallerError(err)
			}
		}
	}
	return nil
}

// ForEachConcurrent concurrently performs a specified 'action' on each element of the sequence.
// If 'ctx' is canceled or 'action' returns non-nil error,
// operation is stopped and corresponding error is returned.
func ForEachConcurrent[T any](ctx context.Context, seq iter.Seq[T], action func(T) error) error {
	if seq == nil {
		return errorhelper.CallerError(ErrNilSec)
	}
	if action == nil {
		return errorhelper.CallerError(ErrNilAction)
	}
	g := new(errgroup.Group)
	for t := range seq {
		g.Go(func() error {
			select {
			case <-ctx.Done():
				return errorhelper.CallerError(ctx.Err())
			default:
				if err := action(t); err != nil {
					return errorhelper.CallerError(err)
				}
			}
			return nil
		})
	}
	return errorhelper.CallerError(g.Wait())
}
