package ratelimiter

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/samber/lo"
)

// Service implements rate limiting per user.
type Service struct {
	mu          sync.RWMutex
	requests    map[int64][]time.Time
	maxRequests int
	window      time.Duration
}

// New creates a new rate limiter.
// maxRequests: maximum number of requests allowed per window
// window: time window for rate limiting
func New(config Config) (*Service, error) {
	if config.MaxRequests <= 0 {
		return nil, fmt.Errorf("%w: MaxRequests must be greater than 0", ErrInvalidConfig)
	}
	if config.Window <= 0 {
		return nil, fmt.Errorf("%w: Window must be greater than 0", ErrInvalidConfig)
	}

	return &Service{
		mu:          sync.RWMutex{},
		requests:    make(map[int64][]time.Time),
		maxRequests: config.MaxRequests,
		window:      config.Window,
	}, nil
}

func (r *Service) Run(ctx context.Context) {
	ticker := time.NewTicker(r.window)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.Cleanup()
		case <-ctx.Done():
			return
		}
	}
}

// Register registers a new request for the given user.
// If the user has exceeded the maximum allowed requests within the window,
// a LimitExceededError is returned.
func (r *Service) Register(userID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	// Clean up old requests outside the window
	validRequests := r.getValidRequests(now, userID)

	// Check if user has exceeded the limit
	if len(validRequests) >= r.maxRequests {
		// Update the requests for this user
		r.requests[userID] = validRequests

		return newLimitExceededError(
			r.maxRequests,
			r.window,
			validRequests[0].Add(r.window),
		)
	}

	// Add current request
	r.requests[userID] = append(r.requests[userID], now)
	return nil
}

// GetRemainingRequests returns the number of remaining requests for a user.
func (r *Service) GetRemainingRequests(userID int64) int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	now := time.Now()

	// Count valid requests within the window
	validRequests := r.getValidRequests(now, userID)

	remaining := r.maxRequests - len(validRequests)
	if remaining < 0 {
		return 0
	}
	return remaining
}

// GetResetTime returns the time when the rate limit resets for a user.
func (r *Service) GetResetTime(userID int64) time.Time {
	r.mu.RLock()
	defer r.mu.RUnlock()

	now := time.Now()

	// Find the oldest request within the window
	validRequests := r.getValidRequests(now, userID)

	if len(validRequests) == 0 {
		return now
	}

	return validRequests[0].Add(r.window)
}

// Cleanup cleans up expired requests.
func (r *Service) Cleanup() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for userID := range r.requests {
		valid := r.getValidRequests(now, userID)
		if len(valid) == 0 {
			delete(r.requests, userID)
		} else {
			r.requests[userID] = valid
		}
	}
}

func (r *Service) getValidRequests(now time.Time, userID int64) []time.Time {
	return lo.DropWhile(
		r.requests[userID],
		func(reqTime time.Time) bool {
			return now.Sub(reqTime) >= r.window
		},
	)
}
