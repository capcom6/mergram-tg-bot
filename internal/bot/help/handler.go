package help

import (
	"fmt"

	"github.com/capcom6/mergram-tg-bot/internal/bot/handler"
	"github.com/capcom6/mergram-tg-bot/pkg/telegofx"
	"github.com/mymmrac/telego"

	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

type Handler struct {
}

func New() handler.Handler {
	return &Handler{}
}

func (h *Handler) Register(router *telegofx.Router) {
	router.HandleMessage(h.Handle, th.CommandEqual("help"))
}

func (h *Handler) Handle(ctx *th.Context, message telego.Message) error {
	// helpText := "Available commands:\n\n" +
	// 	"• /start - Welcome message and bot introduction\n" +
	// 	"• /mermaid - Render Mermaid diagrams\n" +
	// 	"• /help - Show this help message\n\n" +
	// 	"Use /mermaid followed by your diagram code to generate diagrams.\n" +
	// 	"For example: `/mermaid graph TD; A-->B; B-->C;`"
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
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
