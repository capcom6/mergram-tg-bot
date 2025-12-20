package handler

import (
	"context"

	"github.com/capcom6/mergram-tg-bot/pkg/telegofx"
	"github.com/mymmrac/telego"

	tu "github.com/mymmrac/telego/telegoutil"
	"go.uber.org/zap"
)

type Base struct {
	Bot    *telegofx.Bot
	Logger *zap.Logger
}

func (b *Base) HandleError(ctx context.Context, chatID telego.ChatID, err error) {
	b.Logger.Error("handle update failed", zap.Error(err))

	_, sendErr := b.Bot.SendMessage(
		ctx,
		tu.Message(
			chatID,
			"Operation failed. Please try again later or contact support.",
		),
	)
	if sendErr != nil {
		b.Logger.Error("send message failed", zap.Error(sendErr))
	}
}
