package sender

import (
	"fmt"
	"regexp"

	"github.com/NotNotQuinn/go-irc/core"
)

// Records the last message sent from us in every channel
var lastMsgPerChannel = make(map[string]*core.Outgoing)

// Applies all filters to a message
func handleFilterForMessage(msg *core.Outgoing) *core.Outgoing {

	if msg.NoFilter {
		return registerSameMessagAvoidence(msg)
	}

	// Mention users
	if !msg.DM && msg.User.Name != "" {
		msg.Message = fmt.Sprintf("@%s, ", msg.User.Name) + msg.Message
	}

	// Filter out commands
	cond, err := regexp.MatchString("^[\\.\\/]", msg.Message)
	if err != nil {
		core.Errors <- fmt.Errorf("regex for irc command check failed: %w", err)
		// assume match
		cond = true
	}
	if cond {
		msg.Message = ". " + msg.Message
	}

	return registerSameMessagAvoidence(msg)
}

// Checks if sending the same message twice in a row, and appends a character
func registerSameMessagAvoidence(msg *core.Outgoing) *core.Outgoing {
	// Same message avoidence
	if lastMsgPerChannel[msg.Channel] != nil && msg.Message == lastMsgPerChannel[msg.Channel].Message {
		msg.Message += " 󠀀" // invis character
	}

	lastMsgPerChannel[msg.Channel] = msg
	return msg
}
