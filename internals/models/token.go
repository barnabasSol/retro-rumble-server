package models

import "time"

type Token struct {
	PlayerId   string    `json:"id"`
	Expiration time.Time `json:"exp"`
}
