package handlers

import (
	"github.com/NotNotQuinn/go-irc/client"
)

type ClientsHandled struct {
	Twitch bool
}

// Attaches all clients and returns which were successfully handled
func Handle(clients *client.ClientCollection) ClientsHandled {
	var handled ClientsHandled

	if clients.Twitch != nil {
		TwitchAttach(clients.Twitch)
		handled.Twitch = true
	}

	return handled
}
