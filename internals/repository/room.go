package repository

import (
	"github.com/redis/go-redis/v9"
)

type Room struct {
	db *redis.Client
}

func NewRoom(redis *redis.Client) *Room {
	return &Room{
		db: redis,
	}
}

func (r Room) CreateRoom() {

}

func (r Room) DeleteRoom() {
}
