package channels

import (
	"github.com/NotNotQuinn/go-irc/core/command/messages"
)

var (
	MessagesIN  = make(chan *messages.Incoming, 50)
	MessagesOUT = make(chan *messages.Outgoing, 50)
	// Although it doesnt seem like much, it allows for good error loggin later on.
	//
	// Errors should only be passed to this channel if there is no other place, and
	// a panic is not sutible
	Errors = make(chan error, 10)
)
