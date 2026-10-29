package renderer

import "time"

type Config struct {
	// Timeout is the total budget of one render, retries included.
	Timeout time.Duration

	// MaxAttempts is the number of render attempts per request.
	MaxAttempts int

	// Backoff is the delay before the second attempt. Every further attempt
	// doubles it.
	Backoff time.Duration
}
