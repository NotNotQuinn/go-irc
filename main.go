package main

import (
	"twitch-bot/client"
	"twitch-bot/handlers"
)

func main() {
	cc, err := client.GetCollection()
	if err != nil {
		panic(err)
	}

	handlers.Handle(cc)
	cc.JoinAll()

	if Connect(cc) != nil {
		panic(err)
	}
}
