package cmd

import "strings"

var (
	sayCommand *Command = &Command{
		Name:        "say",
		Whitelist:   WL_adminOnly,
		Description: "Sends the text provided to chat.",
		Execution: func(c *Context) (*Return, error) {
			return &Return{
				Success:         true,
				AllowIRCCommand: true,
				Reply:           strings.Join(c.Args, " "),
			}, nil
		},
	}
)
