package ratelimit

import (
	"testing"
	"time"
)

func TestNewLimiter(t *testing.T) {
	// Test that NewLimiter initializes a limiter with the correct requestsPerSecond
	requestsPerSecond := 100
	limiter := NewLimiter(requestsPerSecond)

	if limiter == nil {
		t.Fatal("NewLimiter returned nil")
	}

	if limiter.requestsPerSecond != requestsPerSecond {
		t.Errorf("NewLimiter(%d) created limiter with requestsPerSecond = %d",
			requestsPerSecond, limiter.requestsPerSecond)
	}

	// Check that window is initialized to current time (approximately)
	now := time.Now()
	if limiter.window.Sub(now).Abs() > time.Second {
		t.Errorf("NewLimiter window time is not close to current time: %v vs %v",
			limiter.window, now)
	}

	if limiter.count != 0 {
		t.Errorf("NewLimiter initialized count to %d, expected 0", limiter.count)
	}
}

func TestAllow(t *testing.T) {
	t.Run("AllowsRequestsUpToLimit", func(t *testing.T) {
		// Create a limiter with a small limit for testing
		limit := 5
		limiter := NewLimiter(limit)

		// This should allow 'limit' requests
		for i := 0; i < limit; i++ {
			err := limiter.Allow()
			if err != nil {
				t.Errorf("Allow() returned error on request %d of %d: %v", i+1, limit, err)
			}
		}

		// The next request should be denied
		err := limiter.Allow()
		if err == nil {
			t.Error("Allow() did not return error after exceeding limit")
		}
		if err != ErrRateLimitExceeded {
			t.Errorf("Allow() returned wrong error: got %v, want %v", err, ErrRateLimitExceeded)
		}
	})

	t.Run("ResetsAfterOneSecond", func(t *testing.T) {
		// Create a limiter with a small limit for testing
		limit := 3
		limiter := NewLimiter(limit)

		// Use up all allowed requests
		for i := 0; i < limit; i++ {
			err := limiter.Allow()
			if err != nil {
				t.Errorf("Allow() returned error on request %d of %d: %v", i+1, limit, err)
			}
		}

		// Set window to more than a second ago to simulate time passing
		limiter.window = time.Now().Add(-1100 * time.Millisecond)

		// This should now allow another 'limit' requests
		for i := 0; i < limit; i++ {
			err := limiter.Allow()
			if err != nil {
				t.Errorf("Allow() returned error on request %d of %d after reset: %v",
					i+1, limit, err)
			}
		}

		// The next request should be denied again
		err := limiter.Allow()
		if err == nil {
			t.Error("Allow() did not return error after exceeding limit after reset")
		}
	})
}

func TestAllowConcurrent(t *testing.T) {
	// Test that the rate limiter is thread-safe when used concurrently
	limit := 50
	limiter := NewLimiter(limit)

	// Run 10 goroutines, each making 10 requests
	done := make(chan bool)
	allowed := make(chan bool)

	// Start concurrent goroutines
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 10; j++ {
				err := limiter.Allow()
				if err == nil {
					allowed <- true
				}
			}
			done <- true
		}(i)
	}

	// Count allowed requests
	count := 0
	timeout := time.After(2 * time.Second)
	doneCount := 0

outer:
	for {
		select {
		case <-allowed:
			count++
		case <-done:
			doneCount++
			if doneCount == 10 {
				break outer
			}
		case <-timeout:
			t.Fatal("Test timed out")
			break outer
		}
	}

	// We should have exactly 'limit' allowed requests, and no more
	if count > limit {
		t.Errorf("Rate limiter allowed %d requests, expected max %d", count, limit)
	}

	// Note: we might have fewer than 'limit' allowed requests if goroutines
	// weren't all able to run before others used up the limit
}
