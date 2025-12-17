package renderer

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

type Service struct {
	config Config

	renderer Renderer

	logger *zap.Logger
}

func NewService(config Config, renderer Renderer, logger *zap.Logger) *Service {
	return &Service{
		config: config,

		renderer: renderer,

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
