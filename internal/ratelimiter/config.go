package ratelimiter

import "time"

type Config struct {
	MaxRequests int
	Window      time.Duration
}
