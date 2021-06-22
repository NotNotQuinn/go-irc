package cmd

import "github.com/NotNotQuinn/go-irc/core"

// Context used to invoke commands
type Context struct {
	// Incoming message that this context was created from
	Incoming *core.Incoming
	// Parsed message args
	Args []string
	// Alias/Name used to invoke command
	Invocation string
	// Command that is invoked
	Command *Command
	// User who invoked this command
	User *core.User
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
func (r *Return) ToOutgoing(ctx *Context) *core.Outgoing {
	if r == nil {
		return nil
	}
	return &core.Outgoing{
		Platform: core.Twitch,
		Message:  r.Reply,
		Channel:  ctx.Incoming.Channel,
		User:     ctx.Incoming.User,
		DM:       ctx.Incoming.DMs,
		NoFilter: r.NoFilter,
	}
}
