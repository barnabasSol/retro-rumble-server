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
			evJson, err := json.Marshal(ev)
			if err != nil {
				log.Println(err)
			}
			err = c.conn.WriteMessage(websocket.TextMessage, evJson)
			if err != nil {
				log.Println("failed to send", ev)
			}
		case <-ticker.C:
			println("ping")
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.Leave <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		println("pong")
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
				log.Printf("error player leaving: %v", err)
				c.hub.Leave <- c
			}
			break
		}
		msg = bytes.TrimSpace(bytes.Replace(msg, newline, space, -1))
		var gameEvent event.GameEvent
		err = json.Unmarshal(msg, &gameEvent)
		if err != nil {
			log.Println(err)
		}
		handle, found := c.hub.eventHandlers[gameEvent.Type]
		if !found {
			log.Println("event error", err)
			continue
		}
		//? routing events to a handler to be tested
		handle(gameEvent, c)
	}

}
