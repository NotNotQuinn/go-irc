package command

import (
	"fmt"
	"strings"

	"github.com/NotNotQuinn/go-irc/channels"
	cmd "github.com/NotNotQuinn/go-irc/cmd"
	"github.com/NotNotQuinn/go-irc/config"
	"github.com/NotNotQuinn/go-irc/core/command/messages"
	"github.com/NotNotQuinn/go-irc/core/sender/ratelimiter"
)

func HandleMessage(inMsg *messages.Incoming) error {
	if inMsg == nil {
		return nil
	}
	args, err := prepareMessage(inMsg.Message)
	if err != nil {
		return err
	}
	if len(args) == 0 {
		return nil
	}
	commandName := args[0]
	args = args[1:]

	command := cmd.GetCmd(commandName)
	if command == nil {
		return fmt.Errorf("command not found '%s'", commandName)
	}
	if !ratelimiter.CheckCommand(command, inMsg.Channel, inMsg.User) {
		return nil
	}
	if command.Whitelist != cmd.WL_none {
		switch command.Whitelist {
		case cmd.WL_adminOnly:
			if !config.Public.Users.Admins.Inclues(inMsg.User.Name()) {
				// Ignore
				return nil
			}
		}
	}
	context := &cmd.Context{Incoming: *inMsg, Args: args, Invocation: commandName}
	ratelimiter.IncrementCount(command, inMsg.Channel, inMsg.User)
	responce, err := command.Execution(context)
	if responce != nil {
		channels.MessagesOUT <- responce.ToOutgoing(context)
	}
	return err
}

func prepareMessage(messageText string) ([]string, error) {
	messageText = strings.Trim(messageText, " \t\n󠀀⠀")
	if !strings.HasPrefix(messageText, config.Public.Global.CommandPrefix) {
		return []string{}, nil
	}
	args := strings.Split(strings.TrimPrefix(messageText, config.Public.Global.CommandPrefix), " ")
	return args, nil
}
