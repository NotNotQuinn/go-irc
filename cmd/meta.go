package cmd

import (
	"fmt"
	"time"

	"github.com/NotNotQuinn/go-irc/channels"
	"github.com/NotNotQuinn/go-irc/core/command/messages"
)

var (
	// Command name to pointer to command object
	Commands map[string]*Command = make(map[string]*Command)
	// Command alias to command name
	CommandAliasMap map[string]string = make(map[string]string)
)

const NoCommandName = "__ LHLKJHDLKJHSDLKJSHDuhlkghI&#^GRITK#^RFGKbmf vkyfsmrg"

type Command struct {
	Name        string
	Aliases     []string
	Execution   func(*Context) (*Return, error)
	onLoad      func() error
	Data        DataType
	Whitelist   Whitelist
	Description string
	Cooldown    time.Duration
}

type Context struct {
	Incoming   messages.Incoming
	Args       []string
	Invocation string
}

type Return struct {
	Success bool
	Reply   string
}

type DataType map[string]string

func (r *Return) ToOutgoing(ctx *Context) *messages.Outgoing {
	return &messages.Outgoing{
		Platform:        messages.Twitch,
		Message:         r.Reply,
		Channel:         ctx.Incoming.Channel,
		User:            ctx.Incoming.User,
		DM:              ctx.Incoming.DMs,
		IncomingMessage: &ctx.Incoming,
	}
}

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
		channels.Errors <- fmt.Errorf("command does not have a name\n%+v", cmd)
		cmd.Name = NoCommandName
	}
	if cmd.Description == "" {
		cmd.Description = "(no description)"
	}
	if cmd.Cooldown == 0 {
		cmd.Cooldown = time.Second * 5
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

type Whitelist int

const (
	WL_none      Whitelist = 0
	WL_adminOnly Whitelist = 1
)

func GetCmd(name string) *Command {
	command := Commands[name]
	if command == nil {
		command = Commands[CommandAliasMap[name]]
	}
	return command
}

func LoadAll() {
	pingCommand.load()
	commandCommand.load()
	aboutCommand.load()
	githubCommand.load()
	joinCommand.load()
	gachiCommand.load()
}
