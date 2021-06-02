package client

import (
	"github.com/NotNotQuinn/go-irc/config"

	"github.com/gempir/go-twitch-irc/v2"
)

type ClientCollection struct {
	Twitch *twitch.Client
}

func GetCollection() (cc *ClientCollection, err error) {
	twitch, err := getTwitch()
	if err != nil {
		return nil, err
	}

	return &ClientCollection{twitch}, nil
}

func (cc *ClientCollection) JoinAll() error {
	conf, err := config.GetPublic()
	if err != nil {
		return err
	}
	cc.Twitch.Join(conf.Twitch.Channels...)
	return nil
}

func (cc *ClientCollection) Connect() error {
	return cc.Twitch.Connect()
}
