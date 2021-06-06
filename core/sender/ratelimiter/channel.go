package ratelimiter

import "time"

// Bool channels because they seem like they would use less memory

// Records message rate limits of channels
var channelLimits = make(map[string]chan bool)

// Records whisper rate limits
var whisperLimits = make(chan bool, 1)

// Inits channels with data
func Init() {
	select {
	case whisperLimits <- true:
	default:
	}
}

// Invokes whisper cooldown, waiting if not open
func InvokeWhisper() {
	<-whisperLimits
	go func() {
		time.Sleep(time.Second / 10 * 12)
		whisperLimits <- true
	}()
}

// Invokes message cooldown, waiting if not open
func InvokeMessage(channel string) {
	initChannel(channel)
	// wait until the channels have something to continue
	<-channelLimits[channel]

	go func() {
		time.Sleep(time.Second / 10 * 12)
		channelLimits[channel] <- true
	}()
}

// Initializes channel if not initilized
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
