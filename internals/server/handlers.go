package server

import (
	"github.com/barnabasSol/retro-rumble/internals/event"
)

type Handlers map[string]func(ev event.GameEvent, client *Client)

type WsEventHandler struct {
	handlersMap Handlers
}

func NewWsEventHandler() *WsEventHandler {
	return &WsEventHandler{
		handlersMap: make(Handlers),
	}
}

func (h *WsEventHandler) init() {
	h.handlersMap[""] = func(ev event.GameEvent, client *Client) {}
}
