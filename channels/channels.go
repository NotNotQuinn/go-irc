package channels

import (
	"github.com/NotNotQuinn/go-irc/core/command/messages"
)

var (
	// Messages coming from the platforms (twitch) should go through here to be handled
	MessagesIN = make(chan *messages.Incoming, 50)
	// Messages to be sent out should be sent through here to be properly dispatched
	MessagesOUT = make(chan *messages.Outgoing, 50)
	// Although it doesnt seem like much, it allows for good error logging later on.
	//
	// Errors should only be passed to this channel if there is no other place, and
	// a panic is not suitable
	Errors = make(chan error, 10)
)
