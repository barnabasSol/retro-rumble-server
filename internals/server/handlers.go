package server

import "github.com/barnabasSol/retro-rumble/internals/event"

type HandlerType string
type Handlers map[HandlerType]func(ev event.GameEvent, client *Client)

type WsEventHandler struct {
	handlers Handlers
}

func NewWsEventHandler() *WsEventHandler {
	return &WsEventHandler{
		handlers: make(Handlers),
	}
}

func (h *WsEventHandler) init() {
	h.handlers[""] = func(ev event.GameEvent, client *Client) {}
}
