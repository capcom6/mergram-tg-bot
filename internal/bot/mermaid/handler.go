package mermaid

import (
	"fmt"
	"strings"

	"github.com/capcom6/mergram-tg-bot/internal/bot/handler"
	"github.com/capcom6/mergram-tg-bot/internal/renderer"
	"github.com/capcom6/mergram-tg-bot/pkg/telegofx"
	"github.com/mymmrac/telego"
	"go.uber.org/zap"

	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

type Handler struct {
	handler.Base

	rendererSvc *renderer.Service
}

func New(rendererSvc *renderer.Service, logger *zap.Logger) handler.Handler {
	return &Handler{
		Base: handler.Base{Logger: logger},

		rendererSvc: rendererSvc,
	}
}

func (h *Handler) Register(router *telegofx.Router) {
	router.HandleMessage(h.handle, th.CommandEqual("mermaid"))
}

func (h *Handler) handle(ctx *th.Context, message telego.Message) error {
	diagram := strings.TrimSpace(strings.TrimPrefix(message.Text, "/mermaid"))
	if diagram == "" {
		_, _ = ctx.Bot().
			SendMessage(ctx, tu.Message(message.Chat.ChatID(), "Please provide a Mermaid diagram after /mermaid"))
		return nil
	}

	data, err := h.rendererSvc.Render(ctx, diagram)
	if err != nil {
		return h.HandleError(ctx, message, err) //nolint:wrapcheck //already wrapped
	}

	_, err = ctx.Bot().SendPhoto(
		ctx,
		tu.Photo(
			message.Chat.ChatID(),
			tu.FileFromBytes(data, "diagram.png"),
		),
	)
	if err != nil {
		return fmt.Errorf("send photo: %w", err)
	}

	return nil
}
