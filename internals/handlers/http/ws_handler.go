package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type HttpHandler struct {
	wsUpgr *websocket.Upgrader
}

func NewHttpHandler(wsupg *websocket.Upgrader) HttpHandler {
	return HttpHandler{
		wsUpgr: wsupg,
	}
}

func (h HttpHandler) WsUpgrade(w http.ResponseWriter, r *http.Request) {
	c, err := h.wsUpgr.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("error %s when upgrading connection to websocket", err)
		return
	}
	defer c.Close()
}
