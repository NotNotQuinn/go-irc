package twitchHandler

import (
	"fmt"
	"twitch-bot/core/command"

	"github.com/gempir/go-twitch-irc/v2"
)

func Privmsg(client *twitch.Client, msg twitch.PrivateMessage) {
	fmt.Printf("[%s] %s: %s\n", msg.Channel, msg.User.Name, msg.Message)
	resp, reason, err := command.HandleMessage(msg.Message, msg.User.Name, msg.Channel)
	if err != nil {
		fmt.Println(err)
		return
	}
	if reason != "" {
		return
	}
	client.Say(msg.Channel, resp)
}
