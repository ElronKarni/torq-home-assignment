package ratelimit

import (
	"errors"
	"sync"
	"time"
)

var ErrRateLimitExceeded = errors.New("rate limit exceeded")

// Limiter implements a simple rate limiter
type Limiter struct {
	requestsPerSecond int
	window            time.Time
	count             int
	mu                sync.Mutex
}

// NewLimiter creates a new rate limiter with the specified requests per second
func NewLimiter(requestsPerSecond int) *Limiter {
	return &Limiter{
		requestsPerSecond: requestsPerSecond,
		window:            time.Now(),
	}
}

// Allow checks if a request is allowed under the rate limit
func (l *Limiter) Allow() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()

	// If more than a second has passed since the last window,
	// reset the window and count
	if now.Sub(l.window).Seconds() >= 1 {
		l.window = now
		l.count = 0
	}

	// Check if the current request exceeds the limit
	if l.count >= l.requestsPerSecond {
		return ErrRateLimitExceeded
	}

	// Increment the count and allow the request
	l.count++
	return nil
}
