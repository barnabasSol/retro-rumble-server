package models

import "time"

type Room struct {
	Id   string        `json:"id"`
	Name string        `json:"name"`
	TTL  time.Duration `json:"ttl"`
}
