package internal

import (
	"context"

	"github.com/capcom6/mergram-tg-bot/internal/bot"
	"github.com/capcom6/mergram-tg-bot/internal/config"
	"github.com/capcom6/mergram-tg-bot/internal/ratelimiter"
	"github.com/capcom6/mergram-tg-bot/internal/renderer"
	"github.com/capcom6/mergram-tg-bot/internal/server"
	"github.com/capcom6/mergram-tg-bot/internal/telegram"
	"github.com/capcom6/mergram-tg-bot/pkg/mermaidinkfx"
	"github.com/capcom6/mergram-tg-bot/pkg/pondfx"
	"github.com/go-core-fx/fiberfx"
	"github.com/go-core-fx/healthfx"
	"github.com/go-core-fx/httpfx"
	"github.com/go-core-fx/logger"
	"github.com/go-core-fx/telegofx"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Run(version healthfx.Version) {
	fx.New(
		// CORE MODULES
		logger.Module(),
		logger.WithFxDefaultLogger(),
		// badgerfx.Module(),
		// bunfx.Module(),
		fiberfx.Module(),
		// gocqlfx.Module(),
		// gocqlxfx.Module(),
		// goosefx.Module(),
		// gormfx.Module(),
		healthfx.Module(),
		httpfx.Module(),
		// openaifx.Module(),
		// openrouterfx.Module(),
		// redisfx.Module(),
		// sqlfx.Module(),
		// sqlxfx.Module(),
		telegofx.Module(telegofx.WithRouter(), telegofx.WithHealth()),
		// validatorfx.Module(),
		// watermillfx.Module(),
		//
		// APP MODULES
		config.Module(),
		mermaidinkfx.Module(),
		pondfx.Module(),
		telegram.Module(),
		server.Module(),
		bot.Module(),
		// BUSINESS MODULES
		fx.Supply(version),
		renderer.Module(),
		ratelimiter.Module(),
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
