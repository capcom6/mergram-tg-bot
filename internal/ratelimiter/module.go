package ratelimiter

import (
	"context"

	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Module() fx.Option {
	return fx.Module(
		"ratelimiter",
		logger.WithNamedLogger("ratelimiter"),
		fx.Provide(New),
		fx.Invoke(func(svc *Service, lc fx.Lifecycle, logger *zap.Logger) {
			ctx, cancel := context.WithCancel(context.Background())
			waitCh := make(chan struct{})
			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					go func() {
						defer close(waitCh)
						svc.Run(ctx)
					}()
					logger.Info("ratelimiter started")
					return nil
				},
				OnStop: func(ctx context.Context) error {
					cancel()
					select {
					case <-waitCh:
					case <-ctx.Done():
						logger.Warn("ratelimiter stop timed out")
					}
					logger.Info("ratelimiter stopped")
					return nil
				},
			})
		}),
	)
}
