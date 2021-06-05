package cmd

var (
	githubCommand *Command = &Command{
		Name:    "github",
		Aliases: []string{"gh"},
		Execution: func(c *Context) (*Return, error) {
			return &Return{
				Success: true,
				Reply:   "AlienPls For your \"research purposes\" https://github.com/NotNotQuinn/go-irc",
			}, nil
		},
		Description: "Links to github, where the bot code is hosted.",
	}
)
