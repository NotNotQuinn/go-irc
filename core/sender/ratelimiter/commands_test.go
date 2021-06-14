package ratelimiter

import (
	"os"
	"testing"

	"github.com/NotNotQuinn/go-irc/cmd"
	"github.com/NotNotQuinn/go-irc/config"
	wbUser "github.com/NotNotQuinn/go-irc/core/user"
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
		user    wbUser.IUser
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
		user    wbUser.IUser
	}
	tests := []struct {
		name string
		args args
	}{
		{"New case", args{&cmd.Command{Name: "CMDname"}, "jtv", wbUser.User("yourmom")}},
		{"Command name already initilized", args{&cmd.Command{Name: "CMDname"}, "justinfan123", wbUser.User("yourm0m")}},
		{"Command name and channel already initilized", args{&cmd.Command{Name: "CMDname"}, "justinfan123", wbUser.User("yourmother")}},
		{"Duplicate case", args{&cmd.Command{Name: "CMDname"}, "justinfan123", wbUser.User("yourmother")}},
		{"New case", args{&cmd.Command{Name: "OtherCommand"}, "quinndt", wbUser.User("quinndt")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initCommand(tt.args.command, tt.args.channel, tt.args.user)
			if limits[tt.args.channel] == nil ||
				limits[tt.args.channel][tt.args.command.Name] == nil ||
				limits[tt.args.channel][tt.args.command.Name][tt.args.user.Name()] == nil {

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
		user    wbUser.IUser
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
