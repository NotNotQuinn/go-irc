package main

import (
	"fmt"
	"twitch-bot/client"
	"twitch-bot/core/sender"
	"twitch-bot/handlers"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Print("Starting")
	cc, err := client.GetCollection()
	if err != nil {
		panic(err)
	}

	// Dots to show progress, even though they mostly go all at once
	// its a good measure of startup speed changing over time.
	fmt.Print(".")
	handlers.Handle(cc)

	fmt.Print(".")
	cc.JoinAll()

	fmt.Print(".")
	go sender.HandleAllSends(cc)

	fmt.Print(".")
	err = cc.Connect()
	if err != nil {
		panic(err)
	}
}
