package twitchHandler

import (
	"fmt"

	"github.com/gempir/go-twitch-irc/v2"
)

func Connected(client *twitch.Client) {
	fmt.Println("Connected!")
	client.Say("turtoise", "Hi :)")
}
