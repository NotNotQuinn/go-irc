package client

import (
	"github.com/NotNotQuinn/go-irc/config"

	"github.com/gempir/go-twitch-irc/v2"
)

func getTwitch() (*twitch.Client, error) {
	conf, err := config.GetPrivate()
	if err != nil {
		return nil, err
	}
	return twitch.NewClient(conf.Username, conf.Oauth), nil
}
