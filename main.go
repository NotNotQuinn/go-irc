package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/NotNotQuinn/go-irc/client"
	"github.com/NotNotQuinn/go-irc/cmd"
	"github.com/NotNotQuinn/go-irc/core"
	"github.com/NotNotQuinn/go-irc/core/incoming"
	"github.com/NotNotQuinn/go-irc/core/sender"
	"github.com/NotNotQuinn/go-irc/handlers"
)

// The last time the bot started successfully
var lastSuccessfulRestart time.Time

func main() {
	defer recoverFromDisconnect()
	go handleErrors()
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
	handled := handlers.Handle(cc)
	if !handled.Twitch {
		panic(errors.New("twitch handlers not initilized"))
	}
	go sender.HandleAllSends(cc)

	fmt.Print(".")
	err = cc.JoinAll()
	if err != nil {
		panic(err)
	}

	fmt.Print(".")
	err = cc.Connect()
	if err != nil {
		panic(err)
	}
	lastSuccessfulRestart = time.Now()
}

// Handles all errors
func handleErrors() {
	for {
		// although it doesnt seem like much, it allows for good error logging later on.
		// Errors should only be passed to this stream if there is no other place, and
		// a panic is not sutible
		err := <-core.Errors
		fmt.Printf("Error: %+v\n", err)
	}
}

// Increases as restart attempts increace in count
var restartMult = 1

// Max ammount of time since last restart to accept the restartMult
const maxLastRestart = time.Minute * 5

// Attempts to recover from a disconnect, re-panics other errors
func recoverFromDisconnect() {
	if err := recover(); err != nil {
		s := fmt.Sprint(err)
		if strings.Contains(s, "no such host") && strings.Contains(s, "irc.chat.twitch.tv") {
			if time.Since(lastSuccessfulRestart) > maxLastRestart {
				restartMult = 1
			}
			if !(restartMult >= 32) {
				// will never exceed 32
				// max amount of time waited is 8 mins (15 * 2^5 seconds)
				restartMult *= 2
			}
			sleepTime := time.Second * 15 * time.Duration(restartMult)
			fmt.Println("\nConnection interupted, attempting restart in", sleepTime)
			time.Sleep(sleepTime)
			// I dont know what the best option would be to restart main - this is all I could come up with
			// Stack could overflow
			defer recoverFromDisconnect()
			main()
		}
		panic(fmt.Errorf("%w", err))
	}
}
