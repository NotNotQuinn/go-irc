package twitchHandler

import (
	"github.com/NotNotQuinn/go-irc/channels"
	"github.com/NotNotQuinn/go-irc/core/command/messages"
	"github.com/gempir/go-twitch-irc/v2"
)

func Whisper(msg twitch.WhisperMessage) {
	channels.MessagesIN <- messages.NewIncoming(&msg)
}
