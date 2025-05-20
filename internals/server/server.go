package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type GameServer struct {
	db *redis.Client
	c  *websocket.Conn
}

func NewGameServer(redis *redis.Client) *GameServer {
	return &GameServer{
		db: redis,
		c:  &websocket.Conn{},
	}
}

func (g *GameServer) Start() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

}
