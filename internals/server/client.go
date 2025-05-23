package server

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

	"github.com/barnabasSol/retro-rumble/internals/event"
	"github.com/barnabasSol/retro-rumble/internals/models"
	"github.com/gorilla/websocket"
)

type Client struct {
	egress chan event.GameEvent
	hub    *GameHub
	conn   *websocket.Conn
	player *models.Player
}

func (c *Client) writePump() {

}

func (c *Client) readPump() {
	defer func() {
		c.hub.Left <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Printf("error here: %v", err)
			}
			break
		}
		msg = bytes.TrimSpace(bytes.Replace(msg, newline, space, -1))
		var gameEvent event.GameEvent
		err = json.Unmarshal(msg, &gameEvent)
		if err != nil {
			c.hub.NotifyErr <- event.Error{
				PlayerId: c.player.Id,
				Message:  "server rejected unknown event",
			}
		}
		c.hub.Event <- gameEvent
	}

}
