package renderer

import (
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"renderer",
		logger.WithNamedLogger("renderer"),
		fx.Provide(NewService),
	)
}
