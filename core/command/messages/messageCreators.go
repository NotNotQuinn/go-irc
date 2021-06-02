package messages

import (
	"github.com/gempir/go-twitch-irc/v2"
)

func NewIncoming(msg interface{ GetType() twitch.MessageType }) *Incoming {
	switch v := msg.(type) {
	case *twitch.PrivateMessage:
		return &Incoming{
			Platform: Twitch,
			Channel:  &v.Channel,
			Message:  &v.Message,
			User:     &v.User.Name,
			Raw:      (*twitch.Message)(&msg),
		}
	default:
		return &Incoming{
			Platform: Twitch,
			Raw:      (*twitch.Message)(&msg),
		}
	}
}
