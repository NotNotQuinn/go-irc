package main

import (
	"fmt"

	"github.com/NotNotQuinn/go-irc/channels"
	"github.com/NotNotQuinn/go-irc/client"
	cmd "github.com/NotNotQuinn/go-irc/cmd"
	"github.com/NotNotQuinn/go-irc/core/incoming"
	"github.com/NotNotQuinn/go-irc/core/sender"
	"github.com/NotNotQuinn/go-irc/handlers"
)

func main() {
	go func() {
		for {
			// although it doesnt seem like much, it allows for good error loggin later on.
			// Errors should only be passed to this stream if there is no other place, and
			// a panic is not sutible
			err := <-channels.Errors
			fmt.Printf("Error: %+v\n", err)
		}
	}()
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
