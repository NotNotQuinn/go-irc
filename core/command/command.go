package command

import (
	"strings"

	cmd "github.com/NotNotQuinn/go-irc/cmd"
	"github.com/NotNotQuinn/go-irc/config"
	"github.com/NotNotQuinn/go-irc/core/command/messages"
	"github.com/NotNotQuinn/go-irc/core/sender/ratelimiter"
)

// Handles an incoming message, invoking a command if needed
func HandleMessage(inMsg *messages.Incoming) error {
	cmd, ctx := getCommandAndContext(inMsg)
	if cmd != nil && ctx != nil {
		//channels.MessagesOUT <- responce.ToOutgoing(context)
	}
	return nil
}

func getCommandAndContext(inMsg *messages.Incoming) (*cmd.Command, *cmd.Context) {
	if inMsg == nil {
		return nil, nil
	}
	args := prepareMessage(inMsg.Message)
	if len(args) == 0 {
		return nil, nil
	}
	commandName := args[0]
	args = args[1:]

	command := cmd.GetCmd(commandName)
	if command == nil {
		return nil, nil
	}
	if !ratelimiter.CheckCommand(command, inMsg.Channel, inMsg.User) {
		return nil, nil
	}
	if command.Whitelist != cmd.WL_none {
		switch command.Whitelist {
		case cmd.WL_adminOnly:
			if !config.Public.Users.Admins.Inclues(inMsg.User.Name()) {
				// Ignore
				return nil, nil
			}
		}
	}
	context := &cmd.Context{Incoming: *inMsg, Args: args, Invocation: commandName}
	ratelimiter.InvokeCooldown(command, inMsg.Channel, inMsg.User)
	return command, context
}

// Prepares a message, seperating the arguments
func prepareMessage(messageText string) []string {
	messageText = strings.Trim(messageText, " \t\n󠀀⠀")
	if !strings.HasPrefix(messageText, config.Public.Global.CommandPrefix) {
		return []string{}
	}
	args := strings.Split(strings.TrimPrefix(messageText, config.Public.Global.CommandPrefix), " ")
	return args
}
