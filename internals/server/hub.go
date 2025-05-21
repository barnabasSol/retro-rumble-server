package server

import (
	"github.com/barnabasSol/retro-rumble/internals/event"
	"github.com/barnabasSol/retro-rumble/internals/models"
)

type GameHub struct {
	players map[*models.Player]struct{}
	joined  chan *models.Player
	left    chan *models.Player
	Event   chan *event.GameEvent
}

func NewGameHub() *GameHub {
	return &GameHub{
		players: make(map[*models.Player]struct{}),
		joined:  make(chan *models.Player),
		left:    make(chan *models.Player),
	}
}

func (h *GameHub) Run() {

}
