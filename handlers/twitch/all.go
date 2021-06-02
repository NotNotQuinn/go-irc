package twitchHandler

import "github.com/gempir/go-twitch-irc/v2"

func AttachAll(client *twitch.Client) {
	client.OnPrivateMessage(func(p twitch.PrivateMessage) { Privmsg(client, p) })
	client.OnConnect(Connected)
}
