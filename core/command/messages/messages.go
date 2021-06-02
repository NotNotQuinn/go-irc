package messages

import "github.com/gempir/go-twitch-irc/v2"

type Outgoing struct {
	Platform        PlatformType
	Message         *string
	Channel         *string
	User            *string
	IncomingMessage *Incoming
}

type Incoming struct {
	Platform PlatformType
	Channel  *string
	Message  *string
	User     *string
	Raw      *twitch.Message
}

type PlatformType int

const (
	Unknown PlatformType = -1
	Twitch  PlatformType = 0
)
