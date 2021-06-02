package sender

import (
	"errors"
	"fmt"
	"twitch-bot/channels"
	"twitch-bot/client"
	"twitch-bot/core/command/messages"
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
			// also any logging of messages sent
			cc.Twitch.Say(*msg.Channel, *msg.Message)
		case messages.Unknown:
			fmt.Println(errors.New("platform set to 'unknown' for message, message not sent"))
		default:
			panic(fmt.Errorf("unknown platform!! %d", msg.Platform))
		}
	}
}
