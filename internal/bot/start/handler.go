package start

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
	router.HandleMessage(h.Handle, th.CommandEqual("start"))
}

func (h *Handler) Handle(ctx *th.Context, message telego.Message) error {
	_, err := ctx.Bot().SendMessage(
		ctx,
		tu.Message(
			message.Chat.ChatID(),
			"Welcome to MerGram bot!\n\n"+
				"Use `/help` to get a list of available commands.",
		),
	)
	if err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
