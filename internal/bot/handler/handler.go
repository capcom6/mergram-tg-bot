package handler

import "github.com/go-core-fx/telegofx"

type Handler interface {
	Register(router *telegofx.Router)
}
