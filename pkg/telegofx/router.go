package telegofx

import (
	"context"
	"fmt"

	th "github.com/mymmrac/telego/telegohandler"
)

type Router struct {
	*th.BotHandler
}

func NewRouter(_ Config, bot *Bot) (*Router, error) {
	handler, err := th.NewBotHandler(
		bot.Bot,
		bot.Updates(),
	)
	if err != nil {
		return nil, fmt.Errorf("create handler: %w", err)
	}

	return &Router{
		BotHandler: handler,
	}, nil
}

func (r *Router) Start() error {
	if err := r.BotHandler.Start(); err != nil {
		return fmt.Errorf("start handler: %w", err)
	}

	return nil
}

func (r *Router) Stop(ctx context.Context) error {
	if err := r.BotHandler.StopWithContext(ctx); err != nil {
		return fmt.Errorf("stop handler: %w", err)
	}

	return nil
}
