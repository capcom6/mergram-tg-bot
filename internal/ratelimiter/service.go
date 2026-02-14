package ratelimiter

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/samber/lo"
	"go.uber.org/zap"
)

// Service implements rate limiting per user.
type Service struct {
	mu          sync.RWMutex
	requests    map[int64][]time.Time
	maxRequests int
	window      time.Duration
	logger      *zap.Logger
}

// New creates a new rate limiter.
// maxRequests: maximum number of requests allowed per window
// window: time window for rate limiting
func New(config Config, logger *zap.Logger) (*Service, error) {
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
		logger:      logger,
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

		resetTime := validRequests[0].Add(r.window)
		r.logger.Warn("rate limit exceeded",
			zap.Int64("user_id", userID),
			zap.Int("max_requests", r.maxRequests),
			zap.Int("current_requests", len(validRequests)),
			zap.Time("reset_time", resetTime),
		)

		return newLimitExceededError(
			r.maxRequests,
			r.window,
			resetTime,
		)
	}

	// Add current request
	r.requests[userID] = append(validRequests, now)

	r.logger.Debug("rate limit check passed",
		zap.Int64("user_id", userID),
		zap.Int("current_requests", len(validRequests)+1),
		zap.Int("max_requests", r.maxRequests),
	)

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
	cleanedCount := 0
	activeUsers := 0

	for userID := range r.requests {
		valid := r.getValidRequests(now, userID)
		if len(valid) == 0 {
			delete(r.requests, userID)
			cleanedCount++
		} else {
			r.requests[userID] = valid
			activeUsers++
		}
	}

	if cleanedCount > 0 || activeUsers > 0 {
		r.logger.Debug("rate limiter cleanup completed",
			zap.Int("cleaned_users", cleanedCount),
			zap.Int("active_users", activeUsers),
		)
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
