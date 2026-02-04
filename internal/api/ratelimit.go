package api

import (
	"context"
	"sync"
	"time"
)

const (
	DefaultRequestsPerSecond = 10
	DefaultBurstSize         = 10
)

// TokenBucket implements a token bucket rate limiter.
type TokenBucket struct {
	mu           sync.Mutex
	tokens       float64
	maxTokens    float64
	refillRate   float64
	lastRefill   time.Time
	refillPeriod time.Duration
}

// NewTokenBucket creates a new rate limiter.
func NewTokenBucket(requestsPerSecond, burst int) *TokenBucket {
	return &TokenBucket{
		tokens:       float64(burst),
		maxTokens:    float64(burst),
		refillRate:   float64(requestsPerSecond),
		lastRefill:   time.Now(),
		refillPeriod: time.Second,
	}
}

// NewDefaultTokenBucket creates a rate limiter with default settings.
func NewDefaultTokenBucket() *TokenBucket {
	return NewTokenBucket(DefaultRequestsPerSecond, DefaultBurstSize)
}

// Wait blocks until a token is available or the context is cancelled.
func (tb *TokenBucket) Wait(ctx context.Context) error {
	for {
		tb.mu.Lock()
		tb.refill()

		if tb.tokens >= 1 {
			tb.tokens--
			tb.mu.Unlock()

			return nil
		}

		waitTime := tb.timeUntilToken()
		tb.mu.Unlock()

		timer := time.NewTimer(waitTime)
		select {
		case <-ctx.Done():
			timer.Stop()

			return ctx.Err()
		case <-timer.C:
			// Try again
		}
	}
}

// TryAcquire attempts to acquire a token without blocking.
func (tb *TokenBucket) TryAcquire() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refill()

	if tb.tokens >= 1 {
		tb.tokens--

		return true
	}

	return false
}

func (tb *TokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)

	tokensToAdd := elapsed.Seconds() * tb.refillRate
	tb.tokens = min(tb.maxTokens, tb.tokens+tokensToAdd)
	tb.lastRefill = now
}

func (tb *TokenBucket) timeUntilToken() time.Duration {
	if tb.tokens >= 1 {
		return 0
	}

	tokensNeeded := 1 - tb.tokens

	return time.Duration(tokensNeeded/tb.refillRate*float64(time.Second)) + time.Millisecond
}
