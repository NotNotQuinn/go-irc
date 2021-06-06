package ratelimiter

import (
	"time"

	"github.com/NotNotQuinn/go-irc/cmd"
	wbUser "github.com/NotNotQuinn/go-irc/core/user"
)

// Used to store the state of rate limits with channel, command, and user.
var limits = make(map[string]map[string]map[string]chan bool)

// Check if a combination is currently on cooldown
func CheckCommand(command *cmd.Command, channel string, user wbUser.IUser) bool {
	initCommand(command, channel, user)
	return len(limits[channel][command.Name][user.Name()]) != 0
}

// Ensures a command has a channel set up in the mapping
func initCommand(command *cmd.Command, channel string, user wbUser.IUser) {
	if limits[channel] == nil {
		limits[channel] = make(map[string]map[string]chan bool)
	}
	if limits[channel][command.Name] == nil {
		limits[channel][command.Name] = make(map[string]chan bool)
	}
	if limits[channel][command.Name][user.Name()] == nil {
		limits[channel][command.Name][user.Name()] = make(chan bool, 1)
		// because we just created this, it needs to be filled
		limits[channel][command.Name][user.Name()] <- true
	}
}

// Will invoke the cooldown, waiting if it isnt already open
func InvokeCooldown(command *cmd.Command, channel string, user wbUser.IUser) {
	initCommand(command, channel, user)
	<-limits[channel][command.Name][user.Name()]
	go func() {
		time.Sleep(command.Cooldown)
		limits[channel][command.Name][user.Name()] <- true
	}()
}
