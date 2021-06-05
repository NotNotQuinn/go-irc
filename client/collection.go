package client

import (
	"github.com/NotNotQuinn/go-irc/config"
	"github.com/gempir/go-twitch-irc/v2"
)

type ClientCollection struct {
	Twitch *twitch.Client
}

var Singleton *ClientCollection

func GetCollection() (cc *ClientCollection, err error) {
	if Singleton != nil {
		return Singleton, nil
	}
	twitch, err := getTwitch()
	if err != nil {
		return nil, err
	}

	Singleton = &ClientCollection{twitch}
	return Singleton, nil
}

func (cc *ClientCollection) JoinAll() error {
	cc.Twitch.Join(config.Public.Twitch.Channels...)
	return nil
}

func (cc *ClientCollection) Connect() error {
	return cc.Twitch.Connect()
}
