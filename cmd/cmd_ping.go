package cmd

var pingCommand *Command = &Command{
	Name: "ping",
	Execution: func(ctx *Context) (*Return, error) {
		channel := "whispers"
		if ctx.Incoming.Channel != "" {
			channel = "#" + ctx.Incoming.Channel
		}
		return &Return{
			Success: true,
			Reply:   "Pong! " + ctx.Incoming.User + " in " + channel,
		}, nil
	},
	Description: "Responds with the user and channel.",
}
