package models

import (
	"github.com/barnabasSol/retro-rumble/internals/event"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Player struct {
	Action
	conn *websocket.Conn
	send chan []event.GameEvent
	Id   string  `json:"id"`
	PosX float32 `json:"pos_x"`
	PosY float32 `json:"pos_y"`
}

func NewPlayer(conn *websocket.Conn, pos_x, pos_y float32) *Player {
	return &Player{
		Id:     uuid.NewString(),
		send:   make(chan []event.GameEvent),
		conn:   conn,
		PosX:   pos_x,
		PosY:   pos_y,
		Action: Action{},
	}
}
