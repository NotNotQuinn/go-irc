package ratelimiter

import (
	"os"
	"testing"
	"time"

	"github.com/NotNotQuinn/go-irc/cmd"
	"github.com/NotNotQuinn/go-irc/config"
	"github.com/NotNotQuinn/go-irc/core"
)

func TestMain(m *testing.M) {
	err := config.InitForTests("../../../config")
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestCheckCommand(t *testing.T) {
	type args struct {
		command *cmd.Command
		channel string
		user    *core.User
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckCommand(tt.args.command, tt.args.channel, tt.args.user); got != tt.want {
				t.Errorf("CheckCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIgnoreAllCommandLimits(t *testing.T) {
	type args struct {
		stop <-chan bool
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			IgnoreAllCommandLimits(tt.args.stop)
		})
	}
}

func Test_initCommand(t *testing.T) {
	type args struct {
		command *cmd.Command
		channel string
		user    *core.User
	}
	tests := []struct {
		name string
		args args
	}{
		{"New case", args{&cmd.Command{Name: "CMDname"}, "jtv", &core.User{
			ID:        100000,
			Name:      "yourmom",
			TwitchID:  6660,
			FirstSeen: time.Date(1700, 1, 0, 1, 0, 0, 0, time.UTC).Format("2006-01-02 03:04:05"),
		}}},
		{"Command name already initilized", args{&cmd.Command{Name: "CMDname"}, "justinfan123", &core.User{
			ID:        160000,
			Name:      "yourm0m",
			TwitchID:  666000,
			FirstSeen: time.Date(1400, 1, 0, 1, 0, 0, 0, time.UTC).Format("2006-01-02 03:04:05"),
		}}},
		{"Command name and channel already initilized", args{&cmd.Command{Name: "CMDname"}, "justinfan123", &core.User{
			ID:        4999,
			Name:      "yourmother",
			TwitchID:  9994,
			FirstSeen: time.Date(1499, 1, 0, 0, 1, 0, 0, time.UTC).Format("2006-01-02 03:04:05"),
		}}},
		{"Duplicate case", args{&cmd.Command{Name: "CMDname"}, "justinfan123", &core.User{
			ID:        4999,
			Name:      "yourmother",
			TwitchID:  9994,
			FirstSeen: time.Date(1499, 1, 0, 0, 1, 0, 0, time.UTC).Format("2006-01-02 03:04:05"),
		}}},
		{"New case", args{&cmd.Command{Name: "OtherCommand"}, "quinndt", &core.User{
			ID:        1,
			Name:      "quinndt",
			TwitchID:  123123,
			FirstSeen: time.Date(2020, 6, 24, 6, 1, 1, 0, time.UTC).Format("2006-01-02 03:04:05"),
		}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initCommand(tt.args.command, tt.args.channel, tt.args.user)
			if limits[tt.args.channel] == nil ||
				limits[tt.args.channel][tt.args.command.Name] == nil ||
				limits[tt.args.channel][tt.args.command.Name][tt.args.user.Name] == nil {

				t.Errorf("initCommand() did not initilize the mapping.")
			}
		})
	}
}

func Test_initCommandChannelGlobal(t *testing.T) {
	type args struct {
		command *cmd.Command
		channel string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initCommandChannelGlobal(tt.args.command, tt.args.channel)
		})
	}
}

func TestInvokeCooldown(t *testing.T) {
	type args struct {
		command *cmd.Command
		channel string
		user    *core.User
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InvokeCooldown(tt.args.command, tt.args.channel, tt.args.user)
		})
	}
}

func Test_checkCommandChannelGlobal(t *testing.T) {
	type args struct {
		command *cmd.Command
		channel string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkCommandChannelGlobal(tt.args.command, tt.args.channel); got != tt.want {
				t.Errorf("checkCommandChannelGlobal() = %v, want %v", got, tt.want)
			}
		})
	}
}
