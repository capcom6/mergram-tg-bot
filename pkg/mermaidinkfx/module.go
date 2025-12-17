package mermaidinkfx

import (
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"mermaidink",
		logger.WithNamedLogger("mermaidink"),
		fx.Provide(NewClient),
	)
}
