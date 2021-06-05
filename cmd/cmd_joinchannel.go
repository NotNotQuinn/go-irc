package cmd

import (
	"fmt"
	"strings"

	"github.com/NotNotQuinn/go-irc/client"
	"github.com/NotNotQuinn/go-irc/config"
)

var joinCommand *Command = &Command{
	Name:      "joinchannel",
	Whitelist: WL_adminOnly,
	Execution: func(c *Context) (*Return, error) {
		channel := strings.ToLower(c.Args[0])
		if !strings.HasPrefix(channel, "#") {
			return &Return{
				Success: false,
				Reply:   "Channel should be prefixed with '#'.",
			}, nil
		}
		channel = strings.TrimPrefix(channel, "#")
		cc := client.Singleton
		for _, ch := range config.Public.Twitch.Channels {
			if ch == channel {
				return &Return{
					Reply:   fmt.Sprintf("Channel #%s already joined!", ch),
					Success: false,
				}, nil
			}
		}
		config.Public.Twitch.Channels = append(config.Public.Twitch.Channels, channel)
		success, err := config.Public.Save()
		if err != nil {
			return nil, err
		}
		if success {
			cc.Twitch.Join(channel)
			return &Return{
				Success: true,
				Reply:   "Joined #" + channel + ".",
			}, nil
		}
		return &Return{
			Success: false,
			Reply:   "Could not save #" + channel + " to config file, not joined.",
		}, nil
	},
	Description: "Joins a channel perminently.",
}
