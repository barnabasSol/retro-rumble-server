package event

import "encoding/json"

type GameEvent struct {
	Type  string          `json:"type"`
	Event json.RawMessage `json:"data"`
}
