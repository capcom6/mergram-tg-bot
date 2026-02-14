package start

import (
	"fmt"

	"github.com/capcom6/mergram-tg-bot/internal/bot/handler"
	"github.com/capcom6/mergram-tg-bot/pkg/telegofx"
	"github.com/mymmrac/telego"
	"go.uber.org/zap"

	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

type Handler struct {
	handler.Base
}

func New(bot *telegofx.Bot, logger *zap.Logger) handler.Handler {
	return &Handler{
		Base: handler.Base{Logger: logger, Bot: bot},
	}
}

func (h *Handler) Register(router *telegofx.Router) {
	router.HandleMessage(h.Handle, th.CommandEqual("start"))
}

func (h *Handler) Handle(ctx *th.Context, message telego.Message) error {
	h.Logger.Info("start command received",
		zap.Int64("user_id", message.From.ID),
		zap.Int64("chat_id", message.Chat.ChatID().ID),
		zap.Int("message_id", message.MessageID),
	)

	_, err := ctx.Bot().SendMessage(
		ctx,
		tu.Message(
			message.Chat.ChatID(),
			"Welcome to MerGram bot!\n\n"+
				"Use `/help` to get a list of available commands.",
		),
	)
	if err != nil {
		h.Logger.Error("failed to send start message",
			zap.Int64("user_id", message.From.ID),
			zap.Int64("chat_id", message.Chat.ChatID().ID),
			zap.Error(err),
		)
		return fmt.Errorf("send message: %w", err)
	}

	h.Logger.Debug("start message sent successfully",
		zap.Int64("user_id", message.From.ID),
		zap.Int64("chat_id", message.Chat.ChatID().ID),
	)

	return nil
}
