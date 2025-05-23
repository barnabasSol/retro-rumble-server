package server

import (
	"log"
	"net/http"

	"github.com/barnabasSol/retro-rumble/internals/event"
	"github.com/barnabasSol/retro-rumble/internals/models"
)

type GameHub struct {
	clients   map[*Client]struct{}
	Joined    chan *Client
	Left      chan *Client
	NotifyErr chan event.Error
}

func newGameHub() *GameHub {
	return &GameHub{
		clients: map[*Client]struct{}{},
		Joined:  make(chan *Client, 10),
		Left:    make(chan *Client, 10),
	}
}

func (h *GameHub) run() {
	for {
		select {
		case client, ok := <-h.Joined:
			if !ok {
				log.Println("failed to join")
			}
			log.Println(client)
			h.clients[client] = struct{}{}
		case client, ok := <-h.Left:
			if !ok {
				log.Println("left channel closed")
			}
			delete(h.clients, client)
			close(client.egress)
		case err, ok := <-h.NotifyErr:
			if !ok {
				log.Println("left channel closed")
			}
			log.Println(err)
		}
	}
}

func (g *GameHub) serveWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		hub:    g,
		conn:   conn,
		egress: make(chan event.GameEvent, 256),
		player: models.NewPlayer(100, 300),
	}
	client.hub.Joined <- client

	go client.writePump()
	go client.readPump()
}
