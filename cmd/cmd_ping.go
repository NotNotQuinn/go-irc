package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/NotNotQuinn/go-irc/config"
)

var pingCommand *Command = &Command{
	Name: "ping",
	Execution: func(ctx *Context) (*Return, error) {
		var developmentStatus []string
		if os.Getenv("WB_DOCKER") == "true" {
			developmentStatus = append(developmentStatus, "Dockerized")
		}
		if !config.Public.Production {
			developmentStatus = append(developmentStatus, "Development build")
		}

		statusString := fmt.Sprintf("[%s]", strings.Join(developmentStatus, ", "))
		if statusString == "[]" {
			statusString = ""
		}
		var userRepr string
		if ctx.User != nil {
			userRepr = fmt.Sprintf("You are %s, Local ID: %d; Twitch ID: %d; First Seen: %s", ctx.User.Name, ctx.User.ID, ctx.User.TwitchID, ctx.User.FirstSeen)
		}
		return &Return{
			Success: true,
			Reply:   fmt.Sprintf("PONG! %s", strings.Join([]string{statusString, userRepr}, " ")),
		}, nil
	}}
