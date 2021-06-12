package cmd

import (
	"fmt"
	"strings"

	"github.com/NotNotQuinn/go-irc/client"
	"github.com/NotNotQuinn/go-irc/config"
)

var joinCommand *Command = &Command{
	Name:      "joinchannel",
	Aliases:   []string{"partchannel"},
	Whitelist: WL_adminOnly,
	Execution: func(c *Context) (*Return, error) {
		if len(c.Args) == 0 {
			return &Return{
				Reply: "Provide a channel in the form of #<channel> to part/join.",
			}, nil
		}
		var part bool
		if c.Invocation == "partchannel" {
			part = true
		}
		channel := strings.ToLower(c.Args[0])
		if !strings.HasPrefix(channel, "#") {
			return &Return{
				Success: false,
				Reply:   "Channel should be prefixed with '#'.",
			}, nil
		}
		channel = strings.TrimPrefix(channel, "#")
		cc := client.Singleton
		if part {
			if !stringSliceContains(config.Public.Twitch.Channels, channel) {
				return &Return{
					Success: false,
					Reply:   fmt.Sprintf("Channel #%s not joined!", channel),
				}, nil
			}
		} else {
			if stringSliceContains(config.Public.Twitch.Channels, channel) {
				return &Return{
					Success: false,
					Reply:   fmt.Sprintf("Channel #%s already joined!", channel),
				}, nil
			}
		}
		if part {
			var index int
			for i, ch := range config.Public.Twitch.Channels {
				if ch == channel {
					index = i
					break
				}
			}
			config.Public.Twitch.Channels = remove(config.Public.Twitch.Channels, index)
		} else {
			config.Public.Twitch.Channels = append(config.Public.Twitch.Channels, channel)
		}
		_, err := config.Public.Save()
		if err != nil {
			if part {
				return &Return{
					Success: false,
					Reply:   "Could not remove #" + channel + " to config file, not parted.",
				}, err
			} else {
				return &Return{
					Success: false,
					Reply:   "Could not add #" + channel + " to config file, not joined.",
				}, err
			}
		}

		if part {
			cc.Twitch.Depart(channel)
			return &Return{
				Success: true,
				Reply:   "Parted #" + channel + ".",
			}, nil
		} else {
			cc.Twitch.Join(channel)
			return &Return{
				Success: true,
				Reply:   "Joined #" + channel + ".",
			}, nil
		}
	},
	Description: "Joins or parts a channel perminently.",
}

// check if string slice contains an item
func stringSliceContains(s []string, query string) bool {
	for _, item := range s {
		if item == query {
			return true
		}
	}
	return false
}

// Removes an index from a string slice
func remove(slice []string, index int) []string {
	return append(slice[:index], slice[index+1:]...)
}
