package repository

import "github.com/redis/go-redis/v9"

type Player struct {
	db *redis.Client
}

func NewPlayer(redis *redis.Client) *Player {
	return &Player{
		db: redis,
	}
}

func (p *Player) KickOutPlayer() {

}
