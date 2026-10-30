package config

import (
	"fmt"
	"runtime"
	"time"

	"github.com/capcom6/mergram-tg-bot/pkg/mermaidinkfx"
	"github.com/go-core-fx/config"
	"go.uber.org/zap"
)

type openAPIConfig struct {
	Enabled    bool   `koanf:"enabled"`
	PublicHost string `koanf:"public_host"`
	PublicPath string `koanf:"public_path"`
}

type httpConfig struct {
	Address     string        `koanf:"address"`
	ProxyHeader string        `koanf:"proxy_header"`
	Proxies     []string      `koanf:"proxies"`
	OpenAPI     openAPIConfig `koanf:"openapi"`
}

type httpClientConfig struct {
	ProxyURL            string        `koanf:"proxy_url"`
	Bypass              string        `koanf:"bypass"`
	Timeout             time.Duration `koanf:"timeout"`
	MaxIdleConns        int           `koanf:"max_idle_conns"`
	MaxIdleConnsPerHost int           `koanf:"max_idle_conns_per_host"`
	IdleConnTimeout     time.Duration `koanf:"idle_conn_timeout"`
}

type telegramConfig struct {
	Token    string        `koanf:"token"`
	ProxyURL string        `koanf:"proxy_url"`
	Timeout  time.Duration `koanf:"timeout"`
}

type rendererConfig struct {
	Timeout     time.Duration `koanf:"timeout"`
	MaxAttempts int           `koanf:"max_attempts"`
	Backoff     time.Duration `koanf:"backoff"`
}

type mermaidinkConfig struct {
	BaseURL string `koanf:"base_url"`
}

type queueConfig struct {
	MaxConcurrency int `koanf:"max_concurrency"`
}

type ratelimiterConfig struct {
	MaxRequests int           `koanf:"max_requests"`
	Window      time.Duration `koanf:"window"`
}

type Config struct {
	HTTP        httpConfig        `koanf:"http"`
	HTTPClient  httpClientConfig  `koanf:"http_client"`
	Telegram    telegramConfig    `koanf:"telegram"`
	Renderer    rendererConfig    `koanf:"renderer"`
	MermaidInk  mermaidinkConfig  `koanf:"mermaidink"`
	Queue       queueConfig       `koanf:"queue"`
	RateLimiter ratelimiterConfig `koanf:"ratelimiter"`
}

func New(logger *zap.Logger) (Config, error) {
	//nolint:mnd //default values
	cfg := Config{
		HTTP: httpConfig{
			Address:     "127.0.0.1:3000",
			ProxyHeader: "X-Forwarded-For",
			Proxies:     []string{},
			OpenAPI: openAPIConfig{
				Enabled:    true,
				PublicHost: "",
				PublicPath: "",
			},
		},
		HTTPClient: httpClientConfig{
			ProxyURL:            "",
			Bypass:              "",
			Timeout:             0,
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 32,
			IdleConnTimeout:     90 * time.Second,
		},
		Telegram: telegramConfig{
			Token:    "",
			ProxyURL: "",
			Timeout:  time.Minute,
		},
		Renderer: rendererConfig{
			Timeout:     10 * time.Second,
			MaxAttempts: 3,
			Backoff:     500 * time.Millisecond,
		},
		MermaidInk: mermaidinkConfig{
			BaseURL: mermaidinkfx.BaseURL,
		},
		Queue: queueConfig{
			MaxConcurrency: runtime.NumCPU() * 16,
		},
		RateLimiter: ratelimiterConfig{
			MaxRequests: 5,
			Window:      1 * time.Minute,
		},
	}

	if err := config.Load(&cfg); err != nil {
		logger.Error("failed to load config", zap.Error(err))
		return cfg, fmt.Errorf("load config: %w", err)
	}

	logger.Debug("config loaded successfully")

	return cfg, nil
}
