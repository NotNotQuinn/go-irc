package client

import "github.com/gempir/go-twitch-irc/v2"

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

func (cc *ClientCollection) JoinAll() {
	cc.Twitch.Join("quinndt", "turtoise")
}
