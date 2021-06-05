package ratelimiter

import "time"

// Bool channels because they seem like they would use less memory

var channelLimits = make(map[string]chan bool)
var whisperLimits = make(chan bool, 1)

func Init() {
	select {
	case whisperLimits <- true:
	default:
	}
}

func AwaitSendWhisper() {
	<-whisperLimits
	go func() {
		time.Sleep(time.Second / 10 * 12)
		whisperLimits <- true
	}()
}

func AwaitSendMessage(channel string) {
	initChannel(channel)
	// wait until the channels have something to continue
	<-channelLimits[channel]

	go func() {
		time.Sleep(time.Second / 10 * 12)
		channelLimits[channel] <- true
	}()
}

func initChannel(channel string) {
	if channelLimits[channel] == nil {
		// Limit of 1 message at a time, no matter what
		channelLimits[channel] = make(chan bool, 1)
		channelLimits[channel] <- true
	}
}

func CheckChannel(channel string) bool {
	return channelLimits[channel] == nil || len(channelLimits[channel]) != 0
}
