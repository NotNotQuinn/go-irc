package sender

import (
	"github.com/NotNotQuinn/go-irc/core/command/messages"
)

var lastMsgPerChannel = make(map[string]*messages.Outgoing)

func handleFilterForMessage(msg *messages.Outgoing) *messages.Outgoing {
	if lastMsgPerChannel[msg.Channel] != nil && msg.Message == lastMsgPerChannel[msg.Channel].Message {
		msg.Message += " 󠀀" // invis character
	}

	lastMsgPerChannel[msg.Channel] = msg
	return msg
}
