package cmd

import "fmt"

var pingCommand *Command = &Command{
	Name: "ping",
	Execution: func(ctx *Context) (*Return, error) {
		return &Return{
			Success: true,
			Reply:   fmt.Sprintf("PONG! You are %s, Local ID: %d; Twitch ID: %d; First Seen: %s", ctx.User.Name, ctx.User.ID, ctx.User.TwitchID, ctx.User.FirstSeen),
		}, nil
	}}
