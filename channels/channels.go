package channels

import (
	"github.com/NotNotQuinn/go-irc/core/command/messages"
)

var (
	MessagesIN  = make(chan *messages.Incoming, 50)
	MessagesOUT = make(chan *messages.Outgoing, 50)
)
