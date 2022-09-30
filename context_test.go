package cntcontext

import (
	"context"
	"testing"
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
