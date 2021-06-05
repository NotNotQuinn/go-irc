package ratelimiter

import (
	"time"

	"github.com/NotNotQuinn/go-irc/cmd"
)

// Map layer 1 is channel, layer 2 is command, layer 3 is user
var limits = make(map[string]map[string]map[string]chan bool)

func CheckCommand(command *cmd.Command, channel, user string) bool {
	initCommand(command, channel, user)
	return len(limits[channel][command.Name][user]) != 0
}

func initCommand(command *cmd.Command, channel, user string) {
	if limits[channel] == nil {
		limits[channel] = make(map[string]map[string]chan bool)
	}
	if limits[channel][command.Name] == nil {
		limits[channel][command.Name] = make(map[string]chan bool)
	}
	if limits[channel][command.Name][user] == nil {
		limits[channel][command.Name][user] = make(chan bool, 1)
		// because we just created this, it needs to be filled
		limits[channel][command.Name][user] <- true
	}
}

func IncrementCount(command *cmd.Command, channel, user string) {
	initCommand(command, channel, user)
	<-limits[channel][command.Name][user]
	go func() {
		time.Sleep(command.Cooldown)
		limits[channel][command.Name][user] <- true
	}()
}
