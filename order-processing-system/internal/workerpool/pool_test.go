package workerpool

import (
	"context"
	"log/slog"
	"sync/atomic"
	"testing"
	"time"
)

func TestPoolProcessesJobs(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pool := New(2, 4, slog.Default())
	pool.Start(ctx, 2)
	var count int32
	for i := 0; i < 4; i++ {
		if err := pool.Submit(ctx, Job{ID: "x", Handle: func(context.Context) error {
			atomic.AddInt32(&count, 1)
			return nil
		}}); err != nil {
			t.Fatal(err)
		}
	}
	pool.Shutdown()
	if atomic.LoadInt32(&count) != 4 {
		t.Fatalf("expected 4 processed jobs, got %d", count)
	}
	_ = time.Second
}
