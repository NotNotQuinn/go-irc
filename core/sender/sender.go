package sender

import (
	"errors"
	"fmt"

	"github.com/NotNotQuinn/go-irc/channels"
	"github.com/NotNotQuinn/go-irc/client"
	"github.com/NotNotQuinn/go-irc/core/command/messages"
)

func HandleAllSends(cc *client.ClientCollection) {
	for {
		msg := <-channels.MessagesOUT
		if msg == nil {
			continue
		}
		msg = handleFilterForMessage(msg)
		switch msg.Platform {
		case messages.Twitch:
			if msg.DM {
				cc.Twitch.Whisper(*msg.User, *msg.Message)
				continue
			}
			cc.Twitch.Say(*msg.Channel, *msg.Message)
		case messages.Unknown:
			channels.Errors <- errors.New("platform set to 'unknown' for message, message not sent")
		default:
			channels.Errors <- fmt.Errorf("unknown platform!! %d", msg.Platform)
		}
	}
}
