package handler

import (
	"fmt"

	"github.com/mymmrac/telego"

	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"go.uber.org/zap"
)

type Base struct {
	Logger *zap.Logger
}

func (b *Base) HandleError(ctx *th.Context, message telego.Message, err error) error {
	b.Logger.Error("handle update failed", zap.Error(err))

	_, sendErr := ctx.Bot().SendMessage(
		ctx,
		tu.Message(
			message.Chat.ChatID(),
			"Operation failed. Please try again later or contact support.",
		),
	)
	if sendErr != nil {
		return fmt.Errorf("send message: %w", sendErr)
	}

	return nil
}
