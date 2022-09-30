package cntcontext

import (
	"context"
	"sync"
)

type IncrFunc func() uint64

func WithCount(parent context.Context, limit uint64) (context.Context, IncrFunc) {
	ctx, cancel := context.WithCancel(parent)

	var c uint64
	var mu sync.Mutex
	return ctx, func() uint64 {
		mu.Lock()
		defer mu.Unlock()

		c += 1
		if c > limit {
			cancel()
		}

		return c
	}
}
