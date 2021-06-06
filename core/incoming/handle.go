package incoming

import (
	"github.com/NotNotQuinn/go-irc/channels"
	"github.com/NotNotQuinn/go-irc/core/command"
)

// Handles all incoming messages by invoking command handler each time
func HandleAll() {
	for {
		if msg, isMore := <-channels.MessagesIN; isMore {
			if msg != nil {
				if err := command.HandleMessage(msg); err != nil {
					channels.Errors <- err
				}
			}
		} else {
			return
		}
	}
}
