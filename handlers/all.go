package handlers

import (
	"twitch-bot/client"
	twitchHandler "twitch-bot/handlers/twitch"
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
