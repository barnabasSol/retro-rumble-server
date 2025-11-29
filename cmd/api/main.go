package main

import (
	"log"

	"github.com/barnabasSol/retro-rumble/internals/hub"
	"github.com/barnabasSol/retro-rumble/internals/server"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	hub := hub.NewGameHub()
	go hub.Start()
	quic_server := server.NewQuicServer(":9000", hub)
	log.Println("retro rumble is up and running")
	quic_server.Start()

}
