package main

import (
	"fmt"

	"github.com/NotNotQuinn/go-irc/client"
	cmd "github.com/NotNotQuinn/go-irc/cmds"
	"github.com/NotNotQuinn/go-irc/core/incoming"
	"github.com/NotNotQuinn/go-irc/core/sender"
	"github.com/NotNotQuinn/go-irc/errorStream"
	"github.com/NotNotQuinn/go-irc/handlers"
)

func main() {
	go errorStream.Listen()
	go incoming.HandleAll()

	fmt.Print("Starting")
	cmd.LoadAll()

	// Dots to show progress, even though they mostly go all at once
	// its a good measure of startup speed changing over time.
	fmt.Print(".")
	cc, err := client.GetCollection()
	if err != nil {
		panic(err)
	}

	fmt.Print(".")
	handlers.Handle(cc)

	fmt.Print(".")
	cc.JoinAll()
	go sender.HandleAllSends(cc)

	fmt.Print(".")
	err = cc.Connect()
	if err != nil {
		panic(err)
	}
}
