package client

import (
	"twitch-bot/config"

	"github.com/gempir/go-twitch-irc/v2"
)

func getTwitch() (*twitch.Client, error) {
	conf, err := config.Get()
	if err != nil {
		return nil, err
	}
	return twitch.NewClient(conf.Username, conf.Oauth), nil
}
