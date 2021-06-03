package command

import (
	"fmt"
	"strings"

	"github.com/NotNotQuinn/go-irc/channels"
	cmd "github.com/NotNotQuinn/go-irc/cmds"
	"github.com/NotNotQuinn/go-irc/config"
	"github.com/NotNotQuinn/go-irc/core/command/messages"
)

func HandleMessage(inMsg *messages.Incoming) error {
	if inMsg == nil {
		return nil
	}
	args, err := prepareMessage(*inMsg.Message)
	if err != nil {
		return err
	}
	if len(args) == 0 {
		return nil
	}
	commandName := args[0]
	args = args[1:]

	command := cmd.Commands[commandName]
	if command == nil {
		command = cmd.Commands[cmd.CommandAliasMap[commandName]]
	}
	if command == nil {
		return fmt.Errorf("command not found '%s'", commandName)
	}
	context := &cmd.Context{Incoming: *inMsg, Args: args}
	responce, _ := command.Execution(context)
	channels.MessagesOUT <- responce.ToOutgoing(context)
	return nil
}

func prepareMessage(messageText string) ([]string, error) {
	conf, err := config.GetPublic()
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(messageText, conf.Global.CommandPrefix) {
		return nil, nil
	}

	messageText = strings.Trim(messageText, " \t\n󠀀⠀")
	args := strings.Split(strings.TrimPrefix(messageText, conf.Global.CommandPrefix), " ")
	return args, nil
}
