package cmd

var (
	githubCommand = Command{
		Name: "github",
		Execution: func(c *Context) (Return, error) {
			return Return{
				Success: true,
				Reply:   "AlienPls For your \"research purposes\" https://github.com/NotNotQuinn/go-irc",
			}, nil
		},
	}
)
