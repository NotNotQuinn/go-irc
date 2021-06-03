package incoming

import (
	"github.com/NotNotQuinn/go-irc/channels"
	"github.com/NotNotQuinn/go-irc/core/command"
)

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
