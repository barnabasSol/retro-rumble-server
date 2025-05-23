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
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case ev, ok := <-c.egress:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Println(err)
				c.hub.NotifyErr <- event.Error{
					PlayerId: c.player.Id,
					Message:  "couldn't setup erro",
				}
			}
			evJson, err := json.Marshal(ev)
			if err != nil {
				log.Println(err)
				c.hub.NotifyErr <- event.Error{
					PlayerId: c.player.Id,
					Message:  "error marshalling event",
				}
			}
			w.Write(evJson)
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
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
		//routing events to a handler to bed one soon
	}

}
