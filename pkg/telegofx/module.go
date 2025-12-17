package telegofx

import (
	"context"
	"fmt"

	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Module(withRouter bool) fx.Option {
	opts := []fx.Option{
		logger.WithNamedLogger("telegofx"),
		fx.Provide(New),
		fx.Invoke(func(lc fx.Lifecycle, bot *Bot, logger *zap.Logger, sh fx.Shutdowner) {
			ctx, cancel := context.WithCancel(context.Background())
			waitCh := make(chan struct{})

			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					go func() {
						if err := bot.Run(ctx); err != nil {
							logger.Error("bot run failed", zap.Error(err))
							if shErr := sh.Shutdown(); shErr != nil {
								logger.Error("shutdown failed", zap.Error(shErr))
							}
						}
						close(waitCh)
					}()

					logger.Info("bot started")
					return nil
				},
				OnStop: func(ctx context.Context) error {
					cancel()
					select {
					case <-waitCh:
					case <-ctx.Done():
						logger.Warn("bot stop timed out")
						return nil
					}
					logger.Info("bot stopped")
					return nil
				},
			})
		}),
	}
	if withRouter {
		opts = append(opts, RouterModule())
	}

	return fx.Module(
		"telegofx",
		opts...,
	)
}

func RouterModule() fx.Option {
	return fx.Module(
		"router",
		logger.WithNamedLogger("router"),
		fx.Provide(NewRouter),
		fx.Invoke(func(lc fx.Lifecycle, router *Router, logger *zap.Logger, sh fx.Shutdowner) {
			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					go func() {
						if err := router.Start(); err != nil {
							logger.Error("router start failed", zap.Error(err))
							if shErr := sh.Shutdown(); shErr != nil {
								logger.Error("shutdown failed", zap.Error(shErr))
							}
						}
					}()

					logger.Info("router started")
					return nil
				},
				OnStop: func(ctx context.Context) error {
					if err := router.Stop(ctx); err != nil {
						return fmt.Errorf("stop router: %w", err)
					}
					logger.Info("router stopped")
					return nil
				},
			})
		}),
	)
}
