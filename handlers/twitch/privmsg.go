package twitchHandler

import (
	"github.com/NotNotQuinn/go-irc/channels"
	"github.com/NotNotQuinn/go-irc/core/command/messages"
	"github.com/gempir/go-twitch-irc/v2"
)

func Privmsg(msg twitch.PrivateMessage) {
	channels.MessagesIN <- messages.NewIncoming(&msg)
}
