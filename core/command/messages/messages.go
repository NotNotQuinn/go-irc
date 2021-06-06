package messages

import (
	wbUser "github.com/NotNotQuinn/go-irc/core/user"
	"github.com/gempir/go-twitch-irc/v2"
)

type Outgoing struct {
	Platform        PlatformType
	Message         string
	Channel         string
	User            wbUser.User
	DM              bool
	IncomingMessage *Incoming
}

type Incoming struct {
	Platform PlatformType
	Channel  string
	Message  string
	User     wbUser.User
	Raw      *twitch.Message
	DMs      bool
}

type PlatformType int

const (
	Unknown PlatformType = -1
	Twitch  PlatformType = 0
)
