package repository

import (
	"context"
	"encoding/json"

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

func (p *Player) GetPlayer(ctx context.Context, id string) (*models.Player, error) {
	playerJson := p.db.Get(ctx, PlayerKey+id).Val()
	var player models.Player
	err := json.Unmarshal([]byte(playerJson), &player)
	if err != nil {
		return nil, err
	}
	return &player, nil

}

func (p *Player) DeletePlayer(ctx context.Context, player models.Player) error {
	panic("unimplimented")
}
