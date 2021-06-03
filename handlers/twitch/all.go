package twitchHandler

import "github.com/gempir/go-twitch-irc/v2"

func AttachAll(client *twitch.Client) {
	client.OnConnect(Connected)
	client.OnPrivateMessage(Privmsg)
	client.OnWhisperMessage(Whisper)
}
