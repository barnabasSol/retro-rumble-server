package models

import (
	"github.com/google/uuid"
)

type Player struct {
	Id       string  `json:"id"`
	Username string  `json:"username"`
	PosX     float32 `json:"pos_x"`
	PosY     float32 `json:"pos_y"`
}

func NewPlayer(
	pos_x, pos_y float32,
) *Player {
	return &Player{
		Id:   uuid.NewString(),
		PosX: pos_x,
		PosY: pos_y,
	}
}
