package event

import "encoding/json"

type GameEvent struct {
	Type  string          `json:"type"`
	Event json.RawMessage `json:"data"`
}

func NewGameEvent(ev_type string, event json.RawMessage) *GameEvent {
	return &GameEvent{
		Type:  ev_type,
		Event: event,
	}
}
