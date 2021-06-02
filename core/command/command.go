package command

import (
	"strings"
	"twitch-bot/config"
)

func HandleMessage(text, user, channel string) (resp, reason string, err error) {
	args, reason, err := prepareMessage(text)
	if err != nil {
		return "", reason, err
	}
	if len(args) == 0 {
		return "", reason, nil
	}
	command := args[0]
	args = args[1:]

	switch command {
	case "ping":
		resp := "Pong! " + user + " in #" + channel
		return resp, "", nil
	}
	reason = "no-command"
	return "", reason, nil
}

func prepareMessage(msg string) ([]string, string, error) {
	conf, err := config.GetPublic()
	if err != nil {
		return nil, "", err
	}

	if !strings.HasPrefix(msg, conf.Global.CommandPrefix) {
		return nil, "no-prefix", nil
	}

	msg = strings.Trim(msg, " \t\n󠀀⠀")
	args := strings.Split(strings.TrimPrefix(msg, conf.Global.CommandPrefix), " ")
	return args, "", nil
}
