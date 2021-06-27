package core

import (
	"github.com/gempir/go-twitch-irc/v2"
)

// Create an messages.Incoming message from a twitch message
func NewIncoming(msg interface{ GetType() twitch.MessageType }) *Incoming {
	if msg == nil {
		return nil
	}
	switch v := msg.(type) {
	case *twitch.WhisperMessage:
		return &Incoming{
			Platform: Twitch,
			Channel:  "",
			Message:  v.Message,
			User:     AlwaysGetUser(v.User),
			Raw:      (*twitch.Message)(&msg),
			DMs:      true,
		}
	case *twitch.PrivateMessage:
		return &Incoming{
			Platform: Twitch,
			Channel:  v.Channel,
			Message:  v.Message,
			User:     AlwaysGetUser(v.User),
			Raw:      (*twitch.Message)(&msg),
		}
	default:
		return &Incoming{
			Platform: Twitch,
			Raw:      (*twitch.Message)(&msg),
		}
	}
}

// Create an outgoing message from an messages.Incoming message, and a responce
func NewOutgoing(inMsg *Incoming, responce string) *Outgoing {
	if inMsg == nil {
		return nil
	}
	return &Outgoing{
		Platform: inMsg.Platform,
		Message:  responce,
		Channel:  inMsg.Channel,
		User:     inMsg.User,
		DM:       inMsg.DMs,
	}
}

// Create a fake outgoing message
func FakeOutgoing(channel, message string, platform PlatformType) *Outgoing {
	return &Outgoing{
		Platform: platform,
		Message:  message,
		Channel:  channel,
		User:     &User{ID: 0, Name: "", TwitchID: 0, FirstSeen: ""},
	}
}

func FakeIncoming(channel, message string, user *User, DMs bool, platform PlatformType) *Incoming {
	return &Incoming{
		Platform: platform,
		Channel:  channel,
		Message:  message,
		User:     user,
		Raw:      nil,
		DMs:      DMs,
	}
}
