package main

import "twitch-bot/client"

func Connect(clients *client.ClientCollection) error {
	return clients.Twitch.Connect()
}
