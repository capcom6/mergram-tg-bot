package renderer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/alitto/pond/v2"
	"go.uber.org/zap"
)

type Service struct {
	config Config

	renderer Renderer
	queue    pond.Pool

	logger *zap.Logger
}

func NewService(config Config, renderer Renderer, queue pond.Pool, logger *zap.Logger) *Service {
	return &Service{
		config: config,

		renderer: renderer,
		queue:    queue,

		logger: logger,
	}
}

func (s *Service) Render(ctx context.Context, diagram string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, s.config.Timeout)
	defer cancel()

	data, err := s.renderer.Render(ctx, diagram)
	if err != nil {
		return nil, fmt.Errorf("render diagram: %w", err)
	}

	return data, nil
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
