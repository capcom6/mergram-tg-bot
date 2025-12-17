package renderer

import (
	"fmt"

	"github.com/capcom6/mergram-tg-bot/pkg/mermaidinkfx"
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"renderer",
		logger.WithNamedLogger("renderer"),
		fx.Provide(func() (Renderer, error) {
			r, err := mermaidinkfx.NewClient()
			if err != nil {
				return nil, fmt.Errorf("create mermaidink client: %w", err)
			}

			return r, nil
		}),
		fx.Provide(NewService),
	)
}
