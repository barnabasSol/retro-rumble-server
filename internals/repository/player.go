package repository

import (
	"context"

	"github.com/barnabasSol/retro-rumble/internals/models"
	"github.com/redis/go-redis/v9"
)

const PlayerKey = "player:"

type Player struct {
	db *redis.Client
}

func NewPlayer(redis *redis.Client) *Player {
	return &Player{
		db: redis,
	}
}

func (p *Player) AddPlayer(ctx context.Context, player models.Player) error {
	return p.db.Set(ctx, PlayerKey+player.Id, player, 0).Err()
}
