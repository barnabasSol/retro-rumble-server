package event

import (
	"github.com/barnabasSol/retro-rumble/internals/models"
	"github.com/quic-go/quic-go"
)

type EventHandlerMap map[string]func(models.InboundEvent) error

type NewStreamPayload struct {
	Player     *models.Player
	Identifier string
	Stream     *quic.Stream
}
