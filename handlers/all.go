package handlers

import (
	"github.com/NotNotQuinn/go-irc/client"
	twitchHandler "github.com/NotNotQuinn/go-irc/handlers/twitch"
)

type ClientsHandled struct {
	Twitch bool
}

func Handle(clients *client.ClientCollection) ClientsHandled {
	var handled ClientsHandled

	if clients.Twitch != nil {
		twitchHandler.AttachAll(clients.Twitch)
		handled.Twitch = true
	}

	return handled
}
