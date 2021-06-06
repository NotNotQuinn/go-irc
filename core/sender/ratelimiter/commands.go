package ratelimiter

import (
	"time"

	"github.com/NotNotQuinn/go-irc/cmd"
	wbUser "github.com/NotNotQuinn/go-irc/core/user"
)

// Map layer 1 is channel, layer 2 is command, layer 3 is user
var limits = make(map[string]map[string]map[string]chan bool)

func CheckCommand(command *cmd.Command, channel string, user wbUser.IUser) bool {
	initCommand(command, channel, user)
	return len(limits[channel][command.Name][user.Name()]) != 0
}

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

func IncrementCount(command *cmd.Command, channel string, user wbUser.IUser) {
	initCommand(command, channel, user)
	<-limits[channel][command.Name][user.Name()]
	go func() {
		time.Sleep(command.Cooldown)
		limits[channel][command.Name][user.Name()] <- true
	}()
}
