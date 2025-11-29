package models

type Event struct {
	Type  string `json:"type"`
	Event []byte `json:"data"`
}

type InboundEvent struct {
	Ev     *Event
	Player *Player
}
