package api

import (
	"context"
)

var handlers Handlers

type Handler struct {
	Method      string
	Path        string
	HandlerFunc func(ctx *context.Context)
}
type Handlers map[string]Handler

func NewHandlers() {
	handlers = Handlers{}
}

func ChangeHandlerFunc(handler string, handlerFunc func(ctx *context.Context)) {
	if value, ok := handlers[handler]; ok {
		value.HandlerFunc = handlerFunc
	}
}
