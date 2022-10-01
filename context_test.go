package cntcontext

import (
	"context"
	"errors"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestWithCount(t *testing.T) {
	parent := context.Background()

	var limit uint64 = 2
	ctx, incr := WithCount(parent, limit)

	for i := uint64(0); i < limit; i++ {
		want := i + 1

		count := incr()
		if count != want {
			t.Errorf("FromContext == %d want %d", count, want)
		}

		select {
		case err := <-ctx.Done():
			t.Errorf("<-ctx.Done() == %v", err)
		default:
		}
	}

	count := incr()
	if count != 3 {
		t.Errorf("FromContext == %d want 3", count)
	}

	select {
	case <-ctx.Done():
	default:
		t.Error("<-ctx.Done() blocked, but shouldn't have")
	}
}

func TestWithCount_Parallel(t *testing.T) {
	parent := context.Background()

	var limit uint64 = 50
	ctx, incr := WithCount(parent, limit)

	eg := new(errgroup.Group)
	for i := 0; i < 5; i++ {
		eg.Go(func() error {
			for j := 0; j < 10; j++ {
				_ = incr()

				select {
				case <-ctx.Done():
					return errors.New("context done")
				default:
				}
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		t.Errorf("eg.Wait() == %v", err)
	}

	count := incr()
	if count != 51 {
		t.Errorf("FromContext == %d want 101", count)
	}

	select {
	case <-ctx.Done():
	default:
		t.Error("<-ctx.Done() blocked, but shouldn't have")
	}
}
