package config

import (
	"github.com/capcom6/mergram-tg-bot/internal/ratelimiter"
	"github.com/capcom6/mergram-tg-bot/internal/renderer"
	"github.com/capcom6/mergram-tg-bot/internal/telegram"
	"github.com/capcom6/mergram-tg-bot/pkg/mermaidinkfx"
	"github.com/capcom6/mergram-tg-bot/pkg/pondfx"
	"github.com/go-core-fx/fiberfx"
	"github.com/go-core-fx/fiberfx/openapi"
	"github.com/go-core-fx/httpfx"
	"github.com/go-core-fx/logger"
	"github.com/go-core-fx/telegofx"
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
			func(cfg Config) openapi.Config {
				return openapi.Config{
					Enabled:    cfg.HTTP.OpenAPI.Enabled,
					PublicHost: cfg.HTTP.OpenAPI.PublicHost,
					PublicPath: cfg.HTTP.OpenAPI.PublicPath,
				}
			},
		),
		fx.Provide(
			func(cfg Config) httpfx.Config {
				return httpfx.Config{
					ProxyURL:            cfg.HTTPClient.ProxyURL,
					Bypass:              cfg.HTTPClient.Bypass,
					Timeout:             cfg.HTTPClient.Timeout,
					MaxIdleConns:        cfg.HTTPClient.MaxIdleConns,
					MaxIdleConnsPerHost: cfg.HTTPClient.MaxIdleConnsPerHost,
					IdleConnTimeout:     cfg.HTTPClient.IdleConnTimeout,
					// The upstream API is reached over the public CA.
					TLS: httpfx.TLSConfig{
						RootCAFile:          "",
						RootCAPEM:           "",
						RootCAReplaceSystem: false,
					},
				}
			},
		),
		fx.Provide(
			func(cfg Config) telegram.Config {
				return telegram.Config{
					ProxyURL: cfg.Telegram.ProxyURL,
					Timeout:  cfg.Telegram.Timeout,
				}
			},
			func(cfg Config) telegofx.Config {
				return telegofx.Config{
					Token: cfg.Telegram.Token,
					Mode:  telegofx.ModePolling,
				}
			},
		),
		fx.Provide(
			func(cfg Config) pondfx.Config {
				return pondfx.Config{
					MaxConcurrency: cfg.Queue.MaxConcurrency,
				}
			},
		),
		//
		fx.Provide(
			func(cfg Config) renderer.Config {
				return renderer.Config{
					Timeout:     cfg.Renderer.Timeout,
					MaxAttempts: cfg.Renderer.MaxAttempts,
					Backoff:     cfg.Renderer.Backoff,
				}
			},
		),
		fx.Provide(
			func(cfg Config) mermaidinkfx.Config {
				return mermaidinkfx.Config{
					BaseURL: cfg.MermaidInk.BaseURL,
				}
			},
		),
		fx.Provide(
			func(cfg Config) ratelimiter.Config {
				return ratelimiter.Config{
					MaxRequests: cfg.RateLimiter.MaxRequests,
					Window:      cfg.RateLimiter.Window,
				}
			},
		),
	)
}
