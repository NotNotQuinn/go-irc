package cmd

import (
	"fmt"
	"strings"

	"github.com/NotNotQuinn/go-irc/config"
)

var (
	commandCommand *Command = &Command{
		Name:    "help",
		Aliases: []string{"commands"},
		Execution: func(c *Context) (*Return, error) {
			if c.Invocation == "commands" || len(c.Args) == 0 {
				ret := Return{
					Success: true,
					Reply:   "Commands that are availible: ",
				}
				isAdmin := config.Public.Global.Admin_Username == c.Incoming.User
				i := 0
				for name, cmd := range Commands {
					if cmd == nil {
						continue
					}
					if cmd.Whitelist != WL_none {
						if !isAdmin {
							continue
						}
					}
					ret.Reply += name
					if i < len(Commands)-1 {
						ret.Reply += ", "
					}
					i++
				}
				return &ret, nil
			}
			prefix := config.Public.Global.CommandPrefix
			commandName := c.Args[0]
			command := GetCmd(commandName)
			if command == nil {
				return &Return{
					Success: false,
					Reply:   fmt.Sprintf("Command \"%s\" not found. Do %scommands for a list of commands.", commandName, prefix),
				}, nil
			}
			aliases := []string{}
			for _, alias := range command.Aliases {
				aliases = append(aliases, prefix+alias)
			}
			var whitelisted string
			if command.Whitelist != 0 {
				whitelisted = " (whitelisted)"
			}
			var aliasSuffix string
			if len(aliases) != 0 {
				aliasSuffix = strings.Join(aliases, ", ")
			}
			return &Return{
				Success: true,
				Reply:   fmt.Sprintf("Help for %s%s%s: %s%s", prefix, command.Name, aliasSuffix, command.Description, whitelisted),
			}, nil
		},
		Description: "List all commands, or alternatively shows more detail on a specific command.",
	}
)
