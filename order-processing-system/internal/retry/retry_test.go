package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDoRetriesUntilSuccess(t *testing.T) {
	attempts := 0
	err := Do(context.Background(), 3, time.Millisecond, func(context.Context) error {
		attempts++
		if attempts < 2 {
			return errors.New("temporary")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected success: %v", err)
	}
	if attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts)
	}
}
