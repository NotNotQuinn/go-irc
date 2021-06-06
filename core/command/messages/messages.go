package messages

import (
	wbUser "github.com/NotNotQuinn/go-irc/core/user"
	"github.com/gempir/go-twitch-irc/v2"
)

type Outgoing struct {
	// The platform to send the message on
	Platform PlatformType
	// The message text
	Message string
	// The channel to send the message to
	Channel string
	// The user the message is directed at
	User wbUser.IUser
	// Weather the message should be sent privately
	DM bool
	// The incoming message that invoked this outgoing message
	IncomingMessage *Incoming
}

type Incoming struct {
	// The platform the message was sent on
	Platform PlatformType
	// The channel the message was sent in
	Channel string
	// The message text
	Message string
	// The user who sent the message
	User wbUser.IUser
	// The raw message
	Raw *twitch.Message
	// Whether the message was sent privately
	DMs bool
}

// Platform type to seperate twitch from other platforms in the future
type PlatformType int

const (
	// Platform is unknown - unable to work with data associated
	Unknown PlatformType = -1
	// Twitch platform
	Twitch PlatformType = 0
)
