package retry

import (
	"context"
	"math/rand"
	"time"
)

type PermanentError struct{ Err error }

func (e PermanentError) Error() string { return e.Err.Error() }
func (e PermanentError) Unwrap() error { return e.Err }

func Do(ctx context.Context, attempts int, baseDelay time.Duration, fn func(context.Context) error) error {
	if attempts < 1 {
		attempts = 1
	}
	var err error
	for i := 0; i < attempts; i++ {
		if err = fn(ctx); err == nil {
			return nil
		}
		if _, ok := err.(PermanentError); ok {
			return err
		}
		delay := baseDelay * time.Duration(1<<i)
		delay += time.Duration(rand.Int63n(int64(baseDelay + 1)))
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
	return err
}
