package mermaid

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode/utf16"

	"github.com/capcom6/mergram-tg-bot/internal/bot/handler"
	"github.com/capcom6/mergram-tg-bot/internal/ratelimiter"
	"github.com/capcom6/mergram-tg-bot/internal/renderer"
	"github.com/go-core-fx/telegofx"
	"github.com/mymmrac/telego"
	"github.com/samber/lo"
	"go.uber.org/zap"

	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

type contextKey string

const diagramCodeKey contextKey = "diagramCode"

// renderUnavailableMessage is shown when the renderer could not draw the
// diagram because it was busy or too slow, and another try may succeed.
const renderUnavailableMessage = "⏳ The diagram renderer is busy or timed out. Please try again in a moment."

type Handler struct {
	handler.Base

	rendererSvc  *renderer.Service
	ratelimitSvc *ratelimiter.Service
}

func New(
	bot *telegofx.Bot,
	rendererSvc *renderer.Service,
	ratelimitSvc *ratelimiter.Service,
	logger *zap.Logger,
) handler.Handler {
	return &Handler{
		Base: handler.Base{Logger: logger, Bot: bot},

		rendererSvc:  rendererSvc,
		ratelimitSvc: ratelimitSvc,
	}
}

func (h *Handler) Register(router *telegofx.Router) {
	router.HandleMessage(h.handlePrivateMessage, th.CommandEqual("mermaid"))

	group := router.Group(th.AnyMessageWithText())
	group.Use(func(ctx *th.Context, update telego.Update) error {
		if len(update.Message.Entities) == 0 {
			return ctx.Next(update)
		}

		ent, ok := lo.Find(
			update.Message.Entities,
			func(entity telego.MessageEntity) bool {
				return (entity.Type == "pre" || entity.Type == "code") && entity.Language == "mermaid"
			},
		)

		if !ok {
			return ctx.Next(update)
		}

		utfEncodedString := utf16.Encode([]rune(update.Message.Text))
		if ent.Offset+ent.Length > len(utfEncodedString) {
			h.Logger.Error(
				"invalid entity offset and length",
				zap.Int("offset", ent.Offset),
				zap.Int("length", ent.Length),
			)
			return ctx.Next(update)
		}
		runeString := utf16.Decode(utfEncodedString[ent.Offset : ent.Offset+ent.Length])
		diagramCode := string(runeString)

		ctx = ctx.WithValue(diagramCodeKey, diagramCode)

		return ctx.Next(update)
	})

	group.HandleMessage(h.handleGroupMessage, func(ctx context.Context, _ telego.Update) bool {
		return ctx.Value(diagramCodeKey) != nil
	})
}

func (h *Handler) handlePrivateMessage(ctx *th.Context, message telego.Message) error {
	diagram := strings.TrimSpace(strings.TrimPrefix(message.Text, "/mermaid"))
	if diagram == "" {
		h.Logger.Info("empty mermaid command received",
			zap.Int64("user_id", message.From.ID),
			zap.Int64("chat_id", message.Chat.ChatID().ID),
			zap.Int("message_id", message.MessageID),
		)
		_, _ = ctx.Bot().
			SendMessage(ctx, tu.Message(message.Chat.ChatID(), "Please provide a Mermaid diagram after /mermaid"))
		return nil
	}

	h.Logger.Info("processing mermaid command",
		zap.Int64("user_id", message.From.ID),
		zap.Int64("chat_id", message.Chat.ChatID().ID),
		zap.Int("message_id", message.MessageID),
		zap.Int("diagram_length", len(diagram)),
	)

	h.render(ctx, message.Chat.ChatID(), message.MessageID, diagram)

	return nil
}

func (h *Handler) handleGroupMessage(ctx *th.Context, message telego.Message) error {
	// Extract mermaid code
	diagramCode, ok := ctx.Value(diagramCodeKey).(string)
	if !ok {
		h.Logger.Warn("diagram code not found in context",
			zap.Int64("user_id", message.From.ID),
			zap.Int64("chat_id", message.Chat.ChatID().ID),
			zap.Int("message_id", message.MessageID),
		)
		return nil
	}

	h.Logger.Info("processing group mermaid diagram",
		zap.Int64("user_id", message.From.ID),
		zap.Int64("chat_id", message.Chat.ChatID().ID),
		zap.Int("message_id", message.MessageID),
		zap.Int("diagram_length", len(diagramCode)),
	)

	h.render(ctx, message.Chat.ChatID(), message.MessageID, diagramCode)

	return nil
}

func (h *Handler) render(ctx context.Context, chatID telego.ChatID, sourceMessageID int, diagramCode string) {
	if err := h.ratelimitSvc.Register(chatID.ID); err != nil {
		if limitErr, ok := lo.ErrorsAs[ratelimiter.LimitExceededError](err); ok {
			_, err = h.Bot.SendMessage(ctx, tu.Message(
				chatID,
				fmt.Sprintf(
					"⏰ Rate limit exceeded. Try again in %v.",
					time.Until(limitErr.ResetTime).Round(time.Second),
				),
			))
			if err != nil {
				h.Logger.Error("failed to send rate limit exceeded message", zap.Error(err))
			}
			return
		}

		h.HandleError(ctx, chatID, err)
		return
	}

	// Send processing status message
	processingMsg, err := h.Bot.SendMessage(ctx, tu.Message(
		chatID,
		"🔄 Processing diagram...",
	).WithReplyParameters(&telego.ReplyParameters{
		MessageID:                sourceMessageID,
		ChatID:                   chatID,
		AllowSendingWithoutReply: true,
	}))
	if err != nil {
		h.HandleError(ctx, chatID, err)
		return
	}

	h.rendererSvc.RenderAsync(
		diagramCode,
		func(ctx context.Context, b []byte, err error) {
			if err != nil {
				h.Logger.Error("diagram rendering failed",
					zap.Int64("chat_id", chatID.ID),
					zap.Error(err),
				)

				if h.rendererSvc.IsTransient(err) {
					h.editMessage(ctx, chatID, processingMsg.MessageID, renderUnavailableMessage)
					return
				}

				h.HandleError(ctx, chatID, err)
				return
			}

			h.Logger.Info("diagram rendered successfully",
				zap.Int64("chat_id", chatID.ID),
				zap.Int("message_id", processingMsg.MessageID),
				zap.Int("diagram_size", len(b)),
			)

			_, err = h.Bot.EditMessageMedia(
				ctx,
				tu.EditMessageMedia(
					chatID,
					processingMsg.MessageID,
					tu.MediaPhoto(tu.FileFromBytes(b, imageFileName(b))).
						WithCaption("✅ Diagram rendered successfully!"),
				),
			)
			if err != nil {
				h.Logger.Error("failed to edit message media",
					zap.Int64("chat_id", chatID.ID),
					zap.Int("message_id", processingMsg.MessageID),
					zap.Error(err),
				)
				h.HandleError(ctx, chatID, err)
			}
		},
	)
}

// editMessage replaces the text of a message that was already sent.
func (h *Handler) editMessage(ctx context.Context, chatID telego.ChatID, messageID int, text string) {
	if _, err := h.Bot.EditMessageText(ctx, tu.EditMessageText(chatID, messageID, text)); err != nil {
		h.Logger.Error("failed to edit message text",
			zap.Int64("chat_id", chatID.ID),
			zap.Int("message_id", messageID),
			zap.Error(err),
		)
	}
}

// imageFileName names the uploaded image after its actual content type, the
// renderer answers with JPEG while the extension of the name is what the API
// gets to see.
func imageFileName(data []byte) string {
	switch http.DetectContentType(data) {
	case "image/png":
		return "diagram.png"
	case "image/webp":
		return "diagram.webp"
	default:
		return "diagram.jpg"
	}
}
