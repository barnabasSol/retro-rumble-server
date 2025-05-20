package server

import (
	"net/http"

	handlers "github.com/barnabasSol/retro-rumble/internals/handlers/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type GameServer struct {
	db     *redis.Client
	wsUpgr websocket.Upgrader
}

func NewGameServer(redis *redis.Client) *GameServer {
	return &GameServer{
		db: redis,
	}
}

func (g *GameServer) Start() {
	g.wsUpgr = upgrader
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	h := handlers.NewHttpHandler(&g.wsUpgr)

	r.Post("/ws", h.WsUpgrade)
}
