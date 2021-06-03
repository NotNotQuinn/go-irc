package cmd

var (
	aboutCommand = &Command{
		Name: "about",
		Execution: func(c *Context) (Return, error) {
			return Return{
				Success: true,
				Reply:   "I am an epic chat bot, made by @quinndt in Golang, and running since June 2021. 2020Shred",
			}, nil
		},
	}
)
