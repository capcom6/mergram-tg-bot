package internal

import (
	"context"

	"github.com/capcom6/mergram-tg-bot/internal/bot"
	"github.com/capcom6/mergram-tg-bot/internal/config"
	"github.com/capcom6/mergram-tg-bot/internal/renderer"
	"github.com/capcom6/mergram-tg-bot/internal/server"
	"github.com/capcom6/mergram-tg-bot/pkg/mermaidinkfx"
	"github.com/capcom6/mergram-tg-bot/pkg/telegofx"
	"github.com/go-core-fx/fiberfx"
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Run() {
	fx.New(
		// CORE MODULES
		logger.Module(),
		logger.WithFxDefaultLogger(),
		fiberfx.Module(),
		telegofx.Module(true),
		mermaidinkfx.Module(),
		// APP MODULES
		config.Module(),
		server.Module(),
		bot.Module(),
		// BUSINESS MODULES
		renderer.Module(),
		//
		fx.Invoke(func(lc fx.Lifecycle, logger *zap.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					logger.Info("app started")
					return nil
				},
				OnStop: func(_ context.Context) error {
					logger.Info("app stopped")
					return nil
				},
			})
		}),
	).Run()
}
