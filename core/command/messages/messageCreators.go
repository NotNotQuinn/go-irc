package messages

import (
	wbUser "github.com/NotNotQuinn/go-irc/core/user"
	"github.com/gempir/go-twitch-irc/v2"
)

func NewIncoming(msg interface{ GetType() twitch.MessageType }) *Incoming {
	switch v := msg.(type) {
	case *twitch.WhisperMessage:
		return &Incoming{
			Platform: Twitch,
			Channel:  "",
			Message:  v.Message,
			User:     wbUser.User(v.User.Name),
			Raw:      (*twitch.Message)(&msg),
			DMs:      true,
		}
	case *twitch.PrivateMessage:
		return &Incoming{
			Platform: Twitch,
			Channel:  v.Channel,
			Message:  v.Message,
			User:     wbUser.User(v.User.Name),
			Raw:      (*twitch.Message)(&msg),
		}
	default:
		return &Incoming{
			Platform: Twitch,
			Raw:      (*twitch.Message)(&msg),
		}
	}
}

func NewOutgoing(inMsg *Incoming, responce string) *Outgoing {
	if inMsg == nil {
		return &Outgoing{
			Platform:        Unknown,
			Message:         responce,
			Channel:         "",
			User:            wbUser.User(""),
			IncomingMessage: nil,
			DM:              false,
		}
	}
	return &Outgoing{
		Platform:        inMsg.Platform,
		Message:         responce,
		Channel:         inMsg.Channel,
		User:            inMsg.User,
		IncomingMessage: inMsg,
		DM:              inMsg.DMs,
	}
}

func FakeOutgoing(channel, message string, platform PlatformType) *Outgoing {
	return &Outgoing{
		Platform:        platform,
		Message:         message,
		Channel:         channel,
		User:            "",
		IncomingMessage: nil,
	}
}
