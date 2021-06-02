package command

import "github.com/gempir/go-twitch-irc/v2"

type OutgoingMessageSpec struct {
	Platform        PlatformType
	Message         string
	Channel         string
	IncomingMessage *twitch.Message
}

type PlatformType int

const (
	Unknown PlatformType = -1
	Twitch  PlatformType = 0
)
