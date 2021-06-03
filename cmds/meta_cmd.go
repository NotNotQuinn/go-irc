package cmd

import (
	"errors"
	"fmt"

	"github.com/NotNotQuinn/go-irc/channels"
)

var (
	// Command name to pointer to command object
	Commands map[string]*Command = make(map[string]*Command)
	// Command alias to command name
	CommandAliasMap map[string]string = make(map[string]string)
)

const NoCommandName = "__ LHLKJHDLKJHSDLKJSHDuhlkghI&#^GRITK#^RFGKbmf vkyfsmrg"

// Registers a command to the list of availible commands
func (cmd *Command) register() {
	for _, alias := range cmd.Aliases {
		CommandAliasMap[alias] = cmd.Name
	}
	Commands[cmd.Name] = cmd
}

// Ensures commands have defaults set
func (cmd *Command) ensureDefaults() {
	if cmd.onLoad == nil {
		cmd.onLoad = func() error { return nil }
	}
	if cmd.Data == nil {
		cmd.Data = make(DataType)
	}
	if cmd.Aliases == nil {
		cmd.Aliases = []string{}
	}
	if cmd.Execution == nil {
		cmd.Execution = func(c *Context) (Return, error) {
			return Return{
				Success: false,
				Reply:   "This command has no definition :)",
			}, nil
		}
	}
	if cmd.Name == "" {
		channels.Errors <- errors.New("command does not have a name")
		cmd.Name = NoCommandName
	}
}

// Loads a command and registers
func (cmd *Command) load() {
	cmd.ensureDefaults()
	err := cmd.onLoad()
	if err != nil {
		channels.Errors <- fmt.Errorf("error in %s command onLoad:\n  %w", cmd.Name, err)
	}
	cmd.register()
}
