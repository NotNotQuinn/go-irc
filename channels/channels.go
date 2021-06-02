package channels

import (
	"twitch-bot/core/command/messages"
)

var (
	MessagesIN  = make(chan *messages.Incoming, 50)
	MessagesOUT = make(chan *messages.Outgoing, 50)
)
