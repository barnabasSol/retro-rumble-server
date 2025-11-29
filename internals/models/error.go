package models

type Error struct {
	PlayerId string `json:"player_id"`
	Message  string `json:"message"`
}
