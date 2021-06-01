package twitchHandler

import (
	"fmt"

	"github.com/gempir/go-twitch-irc/v2"
)

func Privmsg(msg twitch.PrivateMessage) {
	fmt.Println(msg.Message)
}
