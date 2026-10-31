package renderer

import (
	"github.com/capcom6/mergram-tg-bot/pkg/mermaidinkfx"
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"renderer",
		logger.WithNamedLogger("renderer"),
		fx.Provide(func(client *mermaidinkfx.Client) Renderer {
			return client
		}),
		fx.Provide(NewService),
	)
}
