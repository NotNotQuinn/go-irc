package cmd

import "github.com/NotNotQuinn/go-irc/core/command/messages"

type Command struct {
	Name      string
	Aliases   []string
	Execution func(*Context) (Return, error)
	onLoad    func() error
	Data      DataType
}

type Context struct {
	Incoming messages.Incoming
	Args     []string
}

type Return struct {
	Success bool
	Reply   string
}

type DataType map[string]string

func (r *Return) ToOutgoing(ctx *Context) *messages.Outgoing {
	return &messages.Outgoing{
		Platform:        messages.Twitch,
		Message:         &r.Reply,
		Channel:         ctx.Incoming.Channel,
		User:            ctx.Incoming.User,
		DM:              ctx.Incoming.DMs,
		IncomingMessage: &ctx.Incoming,
	}
}
