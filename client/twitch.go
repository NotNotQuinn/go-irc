package client

import (
	"github.com/NotNotQuinn/go-irc/config"

	"github.com/gempir/go-twitch-irc/v2"
)

// Gets the twitch client
func getTwitch() (*twitch.Client, error) {
	return twitch.NewClient(config.Private.Username, config.Private.Oauth), nil
}
