package ratelimiter

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrInvalidConfig = errors.New("invalid config")
)

type LimitExceededError struct {
	MaxRequests int
	Window      time.Duration
	ResetTime   time.Time
}

func (e LimitExceededError) Error() string {
	return fmt.Sprintf(
		"rate limit exceeded: max requests: %d, window: %s, reset time: %s",
		e.MaxRequests,
		e.Window,
		e.ResetTime,
	)
}

func newLimitExceededError(maxRequests int, window time.Duration, resetTime time.Time) LimitExceededError {
	return LimitExceededError{
		MaxRequests: maxRequests,
		Window:      window,
		ResetTime:   resetTime,
	}
}
