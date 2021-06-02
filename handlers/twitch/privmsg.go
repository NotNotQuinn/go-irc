package twitchHandler

import (
	"fmt"
	"twitch-bot/core/command"
	"twitch-bot/core/command/messages"

	"github.com/gempir/go-twitch-irc/v2"
)

func Privmsg(client *twitch.Client, msg twitch.PrivateMessage) {
	fmt.Printf("[%s] %s: %s\n", msg.Channel, msg.User.Name, msg.Message)
	message := messages.Incoming{
		Platform: messages.Twitch,
		Channel:  &msg.Channel,
		Message:  &msg.Message,
		User:     &msg.User.Name,
	}
	command.HandleMessage(&message)
}
