package client

import (
	"github.com/NotNotQuinn/go-irc/config"
	"github.com/gempir/go-twitch-irc/v2"
)

type ClientCollection struct {
	Twitch *twitch.Client
}

var Singleton *ClientCollection

// Get the collection of clients
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

// Join all channels that should be joined in all clients
func (cc *ClientCollection) JoinAll() error {
	cc.Twitch.Join(config.Public.Twitch.Channels...)
	return nil
}

// Connect to all clients
func (cc *ClientCollection) Connect() error {
	return cc.Twitch.Connect()
}
