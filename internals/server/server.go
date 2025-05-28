package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/redis/go-redis/v9"
)

type GameServer struct {
	db *redis.Client
}

func NewGameServer(redis *redis.Client) *GameServer {
	return &GameServer{
		db: redis,
	}
}

func (g *GameServer) Start() {
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
	r.Use(middleware.Heartbeat("/health"))

	port := os.Getenv("SERVER_PORT")

	srv := http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	hub := newGameHub(g.db)

	r.Get("/ws", hub.serveWS)
	// r.Post("/reconnect", hub.reconnect)

	go hub.run()
	log.Println("listening on port", port)
	srv.ListenAndServe()
}
