package pondfx

import (
	"context"

	"github.com/alitto/pond/v2"
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"pondfx",
		logger.WithNamedLogger("pondfx"),
		fx.Provide(fx.Annotate(
			New,
			fx.ParamTags("", `optional:"true"`),
		)),
		fx.Invoke(func(p pond.Pool, lc fx.Lifecycle) {
			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					return nil
				},
				OnStop: func(_ context.Context) error {
					p.StopAndWait()
					return nil
				},
			})
		}),
	)
}
