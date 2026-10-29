package renderer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/alitto/pond/v2"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

const (
	// maxAttemptTimeout caps a single render attempt, so that a hanging
	// upstream never consumes the whole render budget.
	maxAttemptTimeout = 5 * time.Second

	// firstAttemptWeight gives the first attempt a double share of the render
	// budget.
	firstAttemptWeight = 2

	// backoffMultiplier doubles the delay between two attempts.
	backoffMultiplier = 2
)

type Service struct {
	config Config

	renderer Renderer
	queue    pond.Pool

	logger *zap.Logger
}

func NewService(config Config, renderer Renderer, queue pond.Pool, logger *zap.Logger) (*Service, error) {
	if config.MaxAttempts < 1 {
		return nil, fmt.Errorf("%w: max attempts must be positive, got %d", ErrInvalidConfig, config.MaxAttempts)
	}

	return &Service{
		config: config,

		renderer: renderer,
		queue:    queue,

		logger: logger,
	}, nil
}

// Render draws the diagram and repeats the attempt while the failure looks
// transient and the render budget still allows another try.
func (s *Service) Render(ctx context.Context, diagram string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, s.config.Timeout)
	defer cancel()

	var data []byte

	attempts, err := lo.AttemptWhile(s.config.MaxAttempts, func(index int) (error, bool) {
		if index > 0 && !sleepCtx(ctx, s.backoff(index-1)) {
			return fmt.Errorf("render budget exhausted: %w", ctx.Err()), false
		}

		attemptCtx, attemptCancel := context.WithTimeout(ctx, s.attemptTimeout(ctx, index))
		defer attemptCancel()

		rendered, err := s.renderer.Render(attemptCtx, diagram)
		if err == nil {
			data = rendered
			return nil, false
		}

		err = fmt.Errorf("render attempt %d: %w", index+1, err)

		if !s.retryable(ctx, err) {
			return err, false
		}

		if index < s.config.MaxAttempts-1 {
			s.logger.Warn("diagram render attempt failed, retrying",
				zap.Int("attempt", index+1),
				zap.Duration("delay", s.backoff(index)),
				zap.Error(err),
			)
		}

		return err, true
	})
	if err != nil {
		return nil, fmt.Errorf("render diagram after %d attempt(s): %w", attempts, err)
	}

	return data, nil
}

// attemptTimeout returns the budget of the attempt at the given index.
//
// The budget left for the render is split between the remaining attempts. The
// first attempt gets a double share, because it pays for the cold Mermaid
// render on the upstream side, while a retry usually hits the upstream cache
// and returns in milliseconds. A single attempt never exceeds
// maxAttemptTimeout, so that a hanging upstream cannot consume the whole
// render budget.
func (s *Service) attemptTimeout(ctx context.Context, index int) time.Duration {
	deadline, ok := ctx.Deadline()
	if !ok {
		return maxAttemptTimeout
	}

	attemptsLeft := s.config.MaxAttempts - index
	share := time.Until(deadline) / time.Duration(attemptsLeft)
	if index == 0 {
		share *= firstAttemptWeight
	}

	if share > maxAttemptTimeout {
		return maxAttemptTimeout
	}

	return share
}

// backoff returns the delay before the attempt following the given index: the
// configured Backoff doubled index times, applied before attempt index+1.
func (s *Service) backoff(index int) time.Duration {
	delay := s.config.Backoff
	for range index {
		delay *= backoffMultiplier
	}

	return delay
}

func (s *Service) RenderAsync(diagram string, callback func(context.Context, []byte, error)) {
	once := sync.Once{}

	task := s.queue.Submit(func() {
		start := time.Now()
		data, err := s.Render(s.queue.Context(), diagram)
		duration := time.Since(start)

		if err != nil {
			s.logger.Error("diagram render failed",
				zap.Duration("duration", duration),
				zap.Int("diagram_length", len(diagram)),
				zap.Error(err),
			)
		} else {
			s.logger.Info("diagram render completed",
				zap.Duration("duration", duration),
				zap.Int("diagram_length", len(diagram)),
				zap.Int("output_size", len(data)),
			)
		}

		once.Do(func() {
			callback(s.queue.Context(), data, err)
		})
	})

	go func() {
		if err := task.Wait(); err != nil {
			s.logger.Error("render task failed",
				zap.Error(err),
				zap.Int("diagram_length", len(diagram)),
			)
			once.Do(func() {
				callback(s.queue.Context(), nil, err)
			})
		}
	}()
}

// IsTransient reports whether the given render failure is worth retrying.
func (s *Service) IsTransient(err error) bool {
	return s.renderer.IsTransient(err)
}

// retryable reports whether the failed attempt is worth repeating.
func (s *Service) retryable(ctx context.Context, err error) bool {
	if ctx.Err() != nil {
		return false
	}

	return s.renderer.IsTransient(err)
}

// sleepCtx waits for the delay and reports whether it completed instead of the
// context being done.
func sleepCtx(ctx context.Context, delay time.Duration) bool {
	if delay <= 0 {
		return ctx.Err() == nil
	}

	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}
