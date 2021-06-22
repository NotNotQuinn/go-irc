package command

import (
	"strings"

	"github.com/NotNotQuinn/go-irc/cmd"
	"github.com/NotNotQuinn/go-irc/config"
	"github.com/NotNotQuinn/go-irc/core"
	"github.com/NotNotQuinn/go-irc/core/sender/ratelimiter"
)

// Handles an incoming message
//
// Will do the following if needed:
//  * Execute commands
//  * Respond to messages
//  * Abort because of ratelimits
func HandleMessage(inMsg *core.Incoming) error {
	ctx := GetContext(inMsg)
	if ctx != nil && ctx.Command != nil {
		// Handle command
		if !ratelimiter.CheckCommand(ctx.Command, inMsg.Channel, inMsg.User) {
			return nil
		}
		if ctx.Command.Whitelist != cmd.WL_none {
			switch ctx.Command.Whitelist {
			case cmd.WL_adminOnly:
				if !config.Public.Users.Admins.Inclues(inMsg.User.Name) {
					// Ignore
					return nil
				}
			}
		}
		ratelimiter.InvokeCooldown(ctx.Command, inMsg.Channel, inMsg.User)
		res, err := ctx.Command.Execution(ctx)
		core.MessagesOUT <- res.ToOutgoing(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetContext(inMsg *core.Incoming) *cmd.Context {
	if inMsg == nil {
		return nil
	}

	isCMD, args := prepareMessage(inMsg.Message)
	commandName := args[0]
	args = args[1:]
	context := &cmd.Context{Incoming: inMsg, Args: args, Invocation: commandName, User: &inMsg.User}
	if isCMD {
		context.Command = cmd.GetCmd(commandName)
	}
	return context
}

// Prepares a message, seperating the arguments
func prepareMessage(messageText string) (isCMD bool, args []string) {
	messageText = strings.Trim(messageText, " \t\n󠀀⠀")
	isCMD = false
	if strings.HasPrefix(messageText, config.Public.Global.CommandPrefix) {
		isCMD = true
	}
	args = strings.Split(strings.TrimPrefix(messageText, config.Public.Global.CommandPrefix), " ")
	return isCMD, args
}
