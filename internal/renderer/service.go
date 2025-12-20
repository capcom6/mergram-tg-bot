package renderer

import (
	"context"
	"fmt"
	"sync"

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
		data, err := s.Render(s.queue.Context(), diagram)
		once.Do(func() {
			callback(s.queue.Context(), data, err)
		})
	})

	go func() {
		if err := task.Wait(); err != nil {
			s.logger.Error("render diagram", zap.Error(err))
			once.Do(func() {
				callback(s.queue.Context(), nil, err)
			})
		}
	}()
}
