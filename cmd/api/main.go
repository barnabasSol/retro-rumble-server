package main

import (
	"log"

	"github.com/barnabasSol/retro-rumble/internals/db"
	"github.com/barnabasSol/retro-rumble/internals/server"
)

func main() {
	redisClient, err := db.NewRedis()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	log.Println("connected to redis")
	gameServer := server.NewGameServer(redisClient)
	gameServer.Start()
}
