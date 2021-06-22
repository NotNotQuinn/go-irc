package ratelimiter

import (
	"time"

	"github.com/NotNotQuinn/go-irc/cmd"
	"github.com/NotNotQuinn/go-irc/core"
)

// Used to store the state of rate limits with channel, command, and user.
var limits = make(map[string]map[string]map[string]chan bool)

// Stores the state of rate limits for channels & commands (no user)
//
// Cooldowns will be shorter in this limiter
var channelCommandLimit = make(map[string]map[string]chan bool)

// Wheather command cooldowns are enabled.
//
// If false, every check to cooldowns involving commands will return true
// Additionally every invocation of a cooldown will return imediately.
var commandCooldownEnabled bool = true

// Check if a combination is currently on cooldown
func CheckCommand(command *cmd.Command, channel string, user core.User) bool {
	if !commandCooldownEnabled {
		return true
	}
	initCommand(command, channel, user)
	return len(limits[channel][command.Name][user.Name]) != 0 && checkCommandChannelGlobal(command, channel)
}

// Ignores all command ratelimits - sending message ratelimits still apply
//
// To stop ignoring ratelimits simply close the `stop` channel
func IgnoreAllCommandLimits(stop <-chan bool) {
	commandCooldownEnabled = false
	<-stop
	commandCooldownEnabled = true
}

// Ensures a command has a channel set up in the mapping
func initCommand(command *cmd.Command, channel string, user core.User) {
	if limits[channel] == nil {
		limits[channel] = make(map[string]map[string]chan bool)
	}
	if limits[channel][command.Name] == nil {
		limits[channel][command.Name] = make(map[string]chan bool)
	}
	if limits[channel][command.Name][user.Name] == nil {
		limits[channel][command.Name][user.Name] = make(chan bool, 1)
		// because we just created this, it needs to be filled
		limits[channel][command.Name][user.Name] <- true
	}
}

// Make sure command channel combo exists
func initCommandChannelGlobal(command *cmd.Command, channel string) {
	if channelCommandLimit[channel] == nil {
		channelCommandLimit[channel] = make(map[string]chan bool)
	}
	if channelCommandLimit[channel][command.Name] == nil {
		channelCommandLimit[channel][command.Name] = make(chan bool, 1)
		// because we just created this, it needs to be filled
		channelCommandLimit[channel][command.Name] <- true
	}
}

// Will invoke the cooldown, waiting if it isnt already open
func InvokeCooldown(command *cmd.Command, channel string, user core.User) {
	if !commandCooldownEnabled {
		return
	}
	initCommand(command, channel, user)
	<-channelCommandLimit[channel][command.Name]
	<-limits[channel][command.Name][user.Name]
	go func() {
		time.Sleep(command.GlobalCooldown)
		channelCommandLimit[channel][command.Name] <- true
	}()
	go func() {
		time.Sleep(command.Cooldown)
		limits[channel][command.Name][user.Name] <- true
	}()
}

// check command channel global cooldown
func checkCommandChannelGlobal(command *cmd.Command, channel string) bool {
	if !commandCooldownEnabled {
		return true
	}
	initCommandChannelGlobal(command, channel)
	return len(channelCommandLimit[channel][command.Name]) != 0
}
