package handlers

import (
	"github.com/NotNotQuinn/go-irc/client"
)

type ClientsHandled struct {
	Twitch bool
}

func Handle(clients *client.ClientCollection) ClientsHandled {
	var handled ClientsHandled

	if clients.Twitch != nil {
		TwitchAttach(clients.Twitch)
		handled.Twitch = true
	}

	return handled
}
