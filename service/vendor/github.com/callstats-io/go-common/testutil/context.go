package testutil

import (
	"context"
	"errors"
	"time"
)

// Errors
var (
	ErrDeadlineExpired = errors.New("Context deadline expired before test completed")
)

// WithCancelContext creates a context, calls the callback with it and cancels it when the function returns
func WithCancelContext(cb func(context.Context)) {
	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()
	cb(ctx)
}

// WithDeadlineContext creates a context with the specified deadline, calls the callback with it and cancels it when the function returns
func WithDeadlineContext(d time.Duration, cb func(context.Context)) error {
	deadline := time.Now().Add(d)
	ctx, cancelCtx := context.WithDeadline(context.Background(), deadline)
	defer cancelCtx()
	cb(ctx)
	if time.Now().After(deadline) {
		return ErrDeadlineExpired
	}
	return nil
}
