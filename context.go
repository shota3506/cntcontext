package cntcontext

import (
	"context"
	"sync"
)

type ctxKey struct{}

type IncrFunc func()

func WithCount(parent context.Context, limit uint64) (context.Context, IncrFunc) {
	ctx, cancel := context.WithCancel(parent)

	var c uint64
	ctx = context.WithValue(ctx, ctxKey{}, &c)

	var mu sync.Mutex
	return ctx, func() {
		mu.Lock()
		defer mu.Unlock()

		c += 1
		if c > limit {
			cancel()
		}
	}
}

func FromContext(ctx context.Context) (uint64, bool) {
	val := ctx.Value(ctxKey{})
	c, ok := val.(*uint64)
	if !ok {
		return 0, false
	}
	return *c, true
}
