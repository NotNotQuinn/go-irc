package cmd

import (
	"fmt"
	"time"

	"github.com/NotNotQuinn/go-irc/core"
)

var (
	// Command name to pointer to command object
	Commands map[string]*Command = make(map[string]*Command)
	// Command alias to command name
	CommandAliasMap map[string]string = make(map[string]string)
)

// Command name of commands that should never be invoked or do not have a name
const NoCommandName = ""

// General structure of all commands
type Command struct {
	// Name of the command.
	Name string
	// Alternate names to refer to the command
	Aliases []string
	// Function called when command executed.
	// Default responds with nothing
	Execution func(*Context) (*Return, error)
	// Function called on load.
	onLoad func() error
	// Runtime data
	Data DataType
	// Whitelist type.
	// Default none (0)
	Whitelist Whitelist
	// Description of command
	Description string
	// User cooldown.
	// Default 5 seconds
	Cooldown time.Duration
	// Global channel cooldown.
	// Default 2 seconds
	GlobalCooldown time.Duration
}

// The data type of commands, abstracted because it may change
type DataType map[string]string

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
		cmd.Execution = func(c *Context) (*Return, error) {
			return &Return{
				Success: false,
			}, nil
		}
	}
	if cmd.Name == "" {
		core.Errors <- fmt.Errorf("command does not have a name\n%+v", cmd)
		cmd.Name = NoCommandName
	}
	if cmd.Description == "" {
		cmd.Description = "(no description)"
	}
	if cmd.Cooldown == 0 {
		cmd.Cooldown = time.Second * 5
	}
	if cmd.GlobalCooldown == 0 {
		cmd.GlobalCooldown = time.Second * 2
	}
}

// Loads a command and registers
func (cmd *Command) Load() {
	cmd.ensureDefaults()
	err := cmd.onLoad()
	if err != nil {
		core.Errors <- fmt.Errorf("error in %s command onLoad:\n  %w", cmd.Name, err)
	}
	cmd.register()
}

// Get a command taking into account its aliases
func GetCmd(name string) *Command {
	command := Commands[name]
	if command == nil {
		command = Commands[CommandAliasMap[name]]
	}
	return command
}

// Loads all commands to be accessed from other places
func LoadAll() {
	pingCommand.Load()
	commandCommand.Load()
	aboutCommand.Load()
	githubCommand.Load()
	joinCommand.Load()
	gachiCommand.Load()
	sayCommand.Load()
}
