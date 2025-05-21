package models

type Action struct {
	IsMovingLeft  *bool `json:"is_moving_left,omitempty"`
	IsMovingRight *bool `json:"is_moving_right,omitempty"`
}
