package handlers

import (
	"fmt"

	"github.com/NotNotQuinn/go-irc/channels"
	"github.com/NotNotQuinn/go-irc/core/command/messages"
	"github.com/gempir/go-twitch-irc/v2"
)

// Ataches all twitch handlers
func TwitchAttach(client *twitch.Client) {
	client.OnConnect(connected)
	client.OnPrivateMessage(privmsg)
	client.OnWhisperMessage(whisper)
}

// Called on every whisper
func whisper(msg twitch.WhisperMessage) {
	channels.MessagesIN <- messages.NewIncoming(&msg)
}

// Called on every privmsg
func privmsg(msg twitch.PrivateMessage) {
	channels.MessagesIN <- messages.NewIncoming(&msg)
}

// Called on connect
func connected() {
	fmt.Println("Connected!")
	channels.MessagesOUT <- messages.FakeOutgoing("turtoise", "Hi :)", messages.Twitch)
}
