package sender

import (
	"errors"
	"fmt"

	"github.com/NotNotQuinn/go-irc/channels"
	"github.com/NotNotQuinn/go-irc/client"
	"github.com/NotNotQuinn/go-irc/core/command/messages"
	"github.com/NotNotQuinn/go-irc/core/sender/ratelimiter"
)

func HandleAllSends(cc *client.ClientCollection) {
	ratelimiter.Init()
	for {
		msg := <-channels.MessagesOUT
		// this must be async, because we wait for ratelimits
		go func() {
			if msg == nil {
				return
			}
			if msg.Message == "" {
				fmt.Printf("%+v\nMessage with no text, not sent.\n", msg)
			}
			msg = handleFilterForMessage(msg)
			switch msg.Platform {
			case messages.Twitch:
				if msg.DM {
					ratelimiter.AwaitSendWhisper()
					cc.Twitch.Whisper(msg.User.Name(), msg.Message)
					return
				}
				ping := ""
				if msg.User.Name() != "" {
					ping = fmt.Sprintf("@%s, ", msg.User)
				}
				ratelimiter.AwaitSendMessage(msg.Channel)
				cc.Twitch.Say(msg.Channel, fmt.Sprintf("%s%s", ping, msg.Message))
			case messages.Unknown:
				channels.Errors <- errors.New("platform set to 'unknown' for message, message not sent")
			default:
				channels.Errors <- fmt.Errorf("unknown platform '%d' (message shown below)\n%+v", msg.Platform, msg)
			}
		}()
	}
}
