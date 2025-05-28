package server

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/barnabasSol/retro-rumble/internals/event"
	"github.com/barnabasSol/retro-rumble/internals/models"
	"github.com/barnabasSol/retro-rumble/internals/repository"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type GameHub struct {
	mu        *sync.RWMutex
	GameEvent chan event.GameEvent
	Reconnect chan struct {
		Id   string
		Conn *websocket.Conn
	}
	playerRepo    *repository.Player
	Join          chan *Client
	Leave         chan *Client
	clients       map[string]*Client
	eventHandlers Handlers
}

func newGameHub(db *redis.Client) *GameHub {
	wsHandler := NewWsEventHandler()
	wsHandler.init()
	return &GameHub{
		mu:            &sync.RWMutex{},
		playerRepo:    repository.NewPlayer(db),
		clients:       map[string]*Client{},
		Join:          make(chan *Client, 20),
		Leave:         make(chan *Client, 50),
		eventHandlers: wsHandler.handlersMap,
	}
}

func (h *GameHub) run() {
	for {
		select {
		case client, ok := <-h.Join:
			if !ok {
				log.Println("failed to join")
			}
			h.mu.Lock()
			h.clients[client.player.Id] = client
			h.mu.Unlock()
		case client, ok := <-h.Leave:
			if !ok {
				log.Println("left channel closed")
			}
			go func() {
				t := time.NewTimer(30 * time.Second)
				select {
				case <-t.C:
					h.mu.Lock()
					delete(h.clients, client.player.Id)
					h.mu.Unlock()
					close(client.egress)
				case newInstance := <-client.hub.Reconnect:
					_ = newInstance
				}
			}()
		}
	}
}
func (h *GameHub) serveWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		conn:   conn,
		hub:    h,
		egress: make(chan event.GameEvent, 256),
		player: models.NewPlayer(100, 100),
	}
	h.Join <- client

	go client.writePump()
	go client.readPump()
}

func (h *GameHub) reconnect(w http.ResponseWriter, r *http.Request) {
	// id := r.URL.Query().Get("id")
	// token := r.URL.Query().Get("token")
	panic("not implemented")
}
