package config

import (
	"fmt"
	"runtime"
	"time"

	"github.com/go-core-fx/config"
	"go.uber.org/zap"
)

type httpConfig struct {
	Address     string   `koanf:"address"`
	ProxyHeader string   `koanf:"proxy_header"`
	Proxies     []string `koanf:"proxies"`
}

type telegramConfig struct {
	Token string `koanf:"token"`
}

type rendererConfig struct {
	Timeout time.Duration `koanf:"timeout"`
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
	Telegram    telegramConfig    `koanf:"telegram"`
	Renderer    rendererConfig    `koanf:"renderer"`
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
		},
		Telegram: telegramConfig{
			Token: "",
		},
		Renderer: rendererConfig{
			Timeout: 3 * time.Second,
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
