package mermaidinkfx

import (
	"context"
	"errors"
	"net"
)

var (
	ErrAPIError            = errors.New("api error")
	ErrUpstreamUnavailable = errors.New("upstream unavailable")
	ErrInvalidConfig       = errors.New("invalid config")
)

// IsTransient reports whether the given failure is worth retrying: an
// upstream timeout, a temporary upstream failure or a throttling response.
//
// A canceled context is never transient, and a rejected diagram
// ([ErrAPIError]) never is either - retrying it would fail the same way.
func (c *Client) IsTransient(err error) bool {
	return isTransient(err)
}

func isTransient(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, ErrUpstreamUnavailable) {
		return true
	}

	if errors.Is(err, context.Canceled) {
		return false
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	if netErr, ok := errors.AsType[net.Error](err); ok {
		return netErr.Timeout()
	}

	return false
}
