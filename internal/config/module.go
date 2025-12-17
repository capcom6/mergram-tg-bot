package config

import (
	"github.com/capcom6/mergram-tg-bot/internal/renderer"
	"github.com/capcom6/mergram-tg-bot/pkg/telegofx"
	"github.com/go-core-fx/fiberfx"
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"config",
		logger.WithNamedLogger("config"),
		fx.Provide(
			New,
			fx.Private,
		),
		fx.Provide(
			func(cfg Config) fiberfx.Config {
				return fiberfx.Config{
					Address:     cfg.HTTP.Address,
					ProxyHeader: cfg.HTTP.ProxyHeader,
					Proxies:     cfg.HTTP.Proxies,
				}
			},
		),
		fx.Provide(
			func(cfg Config) telegofx.Config {
				return telegofx.Config{
					Token: cfg.Telegram.Token,
				}
			},
		),
		//
		fx.Provide(
			func(cfg Config) renderer.Config {
				return renderer.Config{
					Timeout: cfg.Renderer.Timeout,
				}
			},
		),
	)
}
