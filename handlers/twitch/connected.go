package twitchHandler

import (
	"fmt"

	"github.com/NotNotQuinn/go-irc/channels"
	"github.com/NotNotQuinn/go-irc/core/command/messages"
)

func Connected() {
	fmt.Println("Connected!")
	channels.MessagesOUT <- messages.FakeOutgoing("turtoise", "Hi :)", messages.Twitch)
}
