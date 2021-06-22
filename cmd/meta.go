package cmd

import (
	"github.com/NotNotQuinn/go-irc/core/command/messages"
	wbUser "github.com/NotNotQuinn/go-irc/core/user"
)

// Context used to invoke commands
type Context struct {
	// Incoming message that this context was created from
	Incoming *messages.Incoming
	// Parsed message args
	Args []string
	// Alias/Name used to invoke command
	Invocation string
	// Command that is invoked
	Command *Command
	// User who invoked this command
	User *wbUser.User
}

// Command whitelist type
type Whitelist int

const (
	// No whitelist
	WL_none Whitelist = iota
	// Admins only
	WL_adminOnly
)

// Data provided from a command execution
type Return struct {
	Success bool
	// The message to reply with
	Reply string
	// Ignore all filtering
	NoFilter bool
}

// Convert the return data to an outgoing message in a context
func (r *Return) ToOutgoing(ctx *Context) *messages.Outgoing {
	if r == nil {
		return nil
	}
	return &messages.Outgoing{
		Platform: messages.Twitch,
		Message:  r.Reply,
		Channel:  ctx.Incoming.Channel,
		User:     ctx.Incoming.User,
		DM:       ctx.Incoming.DMs,
		NoFilter: r.NoFilter,
	}
}
