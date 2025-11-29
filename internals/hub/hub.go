package hub

import (
	"log"

	"github.com/barnabasSol/retro-rumble/internals/event"
	"github.com/barnabasSol/retro-rumble/internals/models"
)

type GameHub struct {
	players         map[*models.Player]struct{}
	pending_players map[*models.Player]struct{}
	Register        chan *models.Player
	RegisterStream  chan *event.NewStreamPayload
	Unregister      chan *models.Player
	GameEvent       chan models.InboundEvent
	EventHandlers   event.EventHandlerMap
}

func NewGameHub() *GameHub {
	return &GameHub{
		players:         make(map[*models.Player]struct{}),
		pending_players: make(map[*models.Player]struct{}),
		RegisterStream:  make(chan *event.NewStreamPayload, 20),
		Register:        make(chan *models.Player, 100),
		GameEvent:       make(chan models.InboundEvent, 256),
		EventHandlers:   make(event.EventHandlerMap),
	}
}

func (h *GameHub) Start() {
	for {
		select {
		case player := <-h.Register:
			h.players[player] = struct{}{}
		case player := <-h.Unregister:
			delete(h.players, player)
		case new_stream := <-h.RegisterStream:
			new_stream.Player.QuicStreams[new_stream.Identifier] = new_stream.Stream
			h.players[new_stream.Player] = struct{}{}
		case inbound := <-h.GameEvent:
			if handle, found := h.EventHandlers[inbound.Ev.Type]; found {
				handle(inbound)
			} else {
				log.Println("invalid event")
			}
		}
	}
}
