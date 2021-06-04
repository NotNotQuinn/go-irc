package twitchHandler

import (
	"fmt"

	"github.com/NotNotQuinn/go-irc/channels"
	"github.com/NotNotQuinn/go-irc/core/command/messages"
	"github.com/gempir/go-twitch-irc/v2"
)

func AttachAll(client *twitch.Client) {
	client.OnConnect(connected)
	client.OnPrivateMessage(privmsg)
	client.OnWhisperMessage(whisper)
}

func whisper(msg twitch.WhisperMessage) {
	channels.MessagesIN <- messages.NewIncoming(&msg)
}

func privmsg(msg twitch.PrivateMessage) {
	channels.MessagesIN <- messages.NewIncoming(&msg)
}

func connected() {
	fmt.Println("Connected!")
	channels.MessagesOUT <- messages.FakeOutgoing("turtoise", "Hi :)", messages.Twitch)
}
