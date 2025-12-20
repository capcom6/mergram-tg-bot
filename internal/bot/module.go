package bot

import (
	"github.com/capcom6/mergram-tg-bot/internal/bot/handler"
	"github.com/capcom6/mergram-tg-bot/internal/bot/help"
	"github.com/capcom6/mergram-tg-bot/internal/bot/mermaid"
	"github.com/capcom6/mergram-tg-bot/internal/bot/start"
	"github.com/capcom6/mergram-tg-bot/pkg/telegofx"
	"github.com/go-core-fx/logger"
	"github.com/mymmrac/telego"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"bot",
		logger.WithNamedLogger("bot"),
		fx.Provide(func() []telego.BotOption {
			return nil
		}),
		// Provide handlers
		fx.Provide(
			fx.Annotate(start.New, fx.ResultTags(`group:"handlers"`)),
			fx.Annotate(mermaid.New, fx.ResultTags(`group:"handlers"`)),
			fx.Annotate(help.New, fx.ResultTags(`group:"handlers"`)),
		),
		// Register handlers
		fx.Invoke(
			fx.Annotate(
				func(handlers []handler.Handler, r *telegofx.Router) {
					for _, h := range handlers {
						h.Register(r)
					}
				},
				fx.ParamTags(`group:"handlers"`),
			),
		),
	)
}
