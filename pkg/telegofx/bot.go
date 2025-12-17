package telegofx

import (
	"context"
	"fmt"

	"github.com/mymmrac/telego"
	"go.uber.org/zap"
)

type Bot struct {
	*telego.Bot

	updates chan telego.Update

	logger *zap.Logger
}

func New(config Config, options []telego.BotOption, logger *zap.Logger) (*Bot, error) {
	options = append(options, telego.WithLogger(&zapLogger{logger}))

	bot, err := telego.NewBot(config.Token, options...)
	if err != nil {
		return nil, fmt.Errorf("create bot: %w", err)
	}

	return &Bot{
		Bot: bot,

		updates: make(chan telego.Update),

		logger: logger,
	}, nil
}

func (b *Bot) Run(ctx context.Context) error {
	updates, err := b.Bot.UpdatesViaLongPolling(ctx, nil)
	if err != nil {
		return fmt.Errorf("get updates: %w", err)
	}

	for {
		select {
		case update := <-updates:
			select {
			case b.updates <- update:
			case <-ctx.Done():
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (b *Bot) Updates() <-chan telego.Update {
	return b.updates
}
