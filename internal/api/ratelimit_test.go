package api

import (
	"context"
	"testing"
	"time"
)

func TestTokenBucket_TryAcquire(t *testing.T) {
	tb := NewTokenBucket(3, 5)

	// Should be able to acquire burst amount immediately
	for i := 0; i < 5; i++ {
		if !tb.TryAcquire() {
			t.Errorf("TryAcquire() = false, want true for token %d", i+1)
		}
	}

	// Should fail after burst is exhausted
	if tb.TryAcquire() {
		t.Error("TryAcquire() = true after burst exhausted, want false")
	}
}

func TestTokenBucket_Wait(t *testing.T) {
	tb := NewTokenBucket(10, 1) // 10 per second, burst of 1

	ctx := context.Background()

	// First should succeed immediately
	start := time.Now()
	if err := tb.Wait(ctx); err != nil {
		t.Errorf("Wait() error = %v", err)
	}

	elapsed := time.Since(start)
	if elapsed > 10*time.Millisecond {
		t.Errorf("First Wait() took %v, expected immediate", elapsed)
	}

	// Second should wait ~100ms
	start = time.Now()
	if err := tb.Wait(ctx); err != nil {
		t.Errorf("Wait() error = %v", err)
	}

	elapsed = time.Since(start)
	if elapsed < 50*time.Millisecond || elapsed > 200*time.Millisecond {
		t.Errorf("Second Wait() took %v, expected ~100ms", elapsed)
	}
}

func TestTokenBucket_CancelContext(t *testing.T) {
	tb := NewTokenBucket(1, 0) // Very slow, no burst

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := tb.Wait(ctx)
	if err == nil {
		t.Error("Wait() should return error on cancelled context")
	}
}
