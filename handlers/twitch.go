package handlers

import (
	"fmt"

	"github.com/NotNotQuinn/go-irc/config"
	"github.com/NotNotQuinn/go-irc/core"
	"github.com/gempir/go-twitch-irc/v2"
)

// Attaches all twitch handlers
func TwitchAttach(client *twitch.Client) {
	client.OnConnect(connected)
	client.OnPrivateMessage(privmsg)
	client.OnWhisperMessage(whisper)
}

// Called on every whisper
func whisper(msg twitch.WhisperMessage) {
	core.MessagesIN <- core.NewIncoming(&msg)
}

// Called on every privmsg
func privmsg(msg twitch.PrivateMessage) {
	core.MessagesIN <- core.NewIncoming(&msg)
}

// Called on connect
func connected() {
	fmt.Println("Connected!")
	if config.Public.Production {
		core.MessagesOUT <- core.FakeOutgoing("turtoise", "Hi :)", core.Twitch)
	}
}
