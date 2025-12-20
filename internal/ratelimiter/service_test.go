package ratelimiter_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/capcom6/mergram-tg-bot/internal/ratelimiter"
)

func TestRateLimiter(t *testing.T) {
	rl, _ := ratelimiter.New(ratelimiter.Config{MaxRequests: 3, Window: time.Minute})

	userID := int64(12345)

	// Test 1: User under limit (should work)
	for i := 1; i <= 3; i++ {
		if rl.Register(userID) != nil {
			t.Errorf("Request %d should be allowed, but was rate limited", i)
		}
	}

	// Test 2: User over limit (should block)
	if rl.Register(userID) == nil {
		t.Error("Request over limit should be blocked")
	}

	// Test 3: User over limit (should return error)
	if rl.Register(userID) == nil {
		t.Error("Request over limit should be blocked")
	}

	// Test 4: Get remaining requests
	remaining := rl.GetRemainingRequests(userID)
	if remaining != 0 {
		t.Errorf("Expected 0 remaining requests, got %d", remaining)
	}

	// Test 5: Different users (should have separate limits)
	otherUserID := int64(54321)
	if rl.Register(otherUserID) != nil {
		t.Error("Different user should have separate rate limit")
	}

	// Test 6: Reset time should be non-zero for users with requests
	oldestRequest := rl.GetResetTime(userID)
	if oldestRequest.IsZero() {
		t.Error("Reset time should not be zero")
	}
}

func TestRateLimiterConcurrency(t *testing.T) {
	rl, _ := ratelimiter.New(ratelimiter.Config{MaxRequests: 3, Window: time.Minute})
	userID := int64(67890)

	// Simulate concurrent requests
	allowed := atomic.Int32{}
	denied := atomic.Int32{}
	done := make(chan bool, 10)
	for range 10 {
		go func() {
			if rl.Register(userID) == nil {
				allowed.Add(1)
			} else {
				denied.Add(1)
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for range 10 {
		<-done
	}

	// Check final state
	remaining := rl.GetRemainingRequests(userID)
	if remaining != 0 {
		t.Errorf("Expected 0 remaining requests, got %d", remaining)
	}
	if a := allowed.Load(); a != 3 {
		t.Errorf("Expected exactly 3 requests to be allowed, got %d", a)
	}
	if d := denied.Load(); d != 7 {
		t.Errorf("Expected exactly 7 requests to be denied, got %d", d)
	}
}
