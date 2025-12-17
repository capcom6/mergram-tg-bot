package handler

import "github.com/capcom6/mergram-tg-bot/pkg/telegofx"

type Handler interface {
	Register(router *telegofx.Router)
}
