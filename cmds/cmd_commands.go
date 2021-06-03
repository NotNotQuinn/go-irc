package cmd

var (
	commandCommand = &Command{
		Name:    "help",
		Aliases: []string{"commands"},
		Execution: func(c *Context) (Return, error) {
			ret := Return{Reply: "Commands that are availible: "}
			i := 0
			for name := range Commands {
				ret.Reply += name
				if i < len(Commands)-1 {
					ret.Reply += ", "
				}
				i++
			}
			return ret, nil
		},
	}
)
