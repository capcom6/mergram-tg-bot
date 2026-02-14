package help

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
	router.HandleMessage(h.Handle, th.CommandEqual("help"))
}

func (h *Handler) Handle(ctx *th.Context, message telego.Message) error {
	h.Logger.Info("help command received",
		zap.Int64("user_id", message.From.ID),
		zap.Int64("chat_id", message.Chat.ChatID().ID),
		zap.Int("message_id", message.MessageID),
	)

	_, err := ctx.Bot().SendMessage(
		ctx,
		tu.MessageWithEntities(
			message.Chat.ChatID(),
			tu.Entity("Available commands:\n\n"),
			tu.Entity("• /start - Welcome message and bot introduction\n"),
			tu.Entity("• /mermaid - Render Mermaid diagrams\n"),
			tu.Entity("• /help - Show this help message\n\n"),
			tu.Entity("Use /mermaid followed by your diagram code to generate diagrams.\n"),
			tu.Entity("For example: "), tu.Entity("/mermaid graph TD; A-->B; B-->C;").Code(),
		),
	)
	if err != nil {
		h.Logger.Error("failed to send help message",
			zap.Int64("user_id", message.From.ID),
			zap.Int64("chat_id", message.Chat.ChatID().ID),
			zap.Error(err),
		)
		return fmt.Errorf("send message: %w", err)
	}

	h.Logger.Debug("help message sent successfully",
		zap.Int64("user_id", message.From.ID),
		zap.Int64("chat_id", message.Chat.ChatID().ID),
	)

	return nil
}
