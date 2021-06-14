package command

import (
	"os"
	"reflect"
	"testing"
	"time"

	cmd "github.com/NotNotQuinn/go-irc/cmd"
	"github.com/NotNotQuinn/go-irc/config"
	"github.com/NotNotQuinn/go-irc/core/command/messages"
	"github.com/NotNotQuinn/go-irc/core/sender/ratelimiter"
)

func TestMain(m *testing.M) {
	err := config.InitForTests("../../config")
	if err != nil {
		panic(err)
	}
	// ignore rate limits
	ignoreLimits := make(chan bool)
	go ratelimiter.IgnoreAllCommandLimits(ignoreLimits)
	code := m.Run()
	close(ignoreLimits)
	os.Exit(code)
}

func TestHandleMessage(t *testing.T) {
	type args struct {
		inMsg *messages.Incoming
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := HandleMessage(tt.args.inMsg); (err != nil) != tt.wantErr {
				t.Errorf("HandleMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getCommandAndContext(t *testing.T) {
	testCommand := &cmd.Command{
		Execution: func(c *cmd.Context) (*cmd.Return, error) {
			return nil, nil
		},
		Name:           "testCMD",
		Aliases:        []string{"AYYYYyyyyyyyLMAAAAAOOOOOO_Alien_Please"},
		Data:           make(cmd.DataType),
		Whitelist:      cmd.WL_none,
		Description:    "Ya yayayaya forsenCD",
		Cooldown:       time.Second * 8,
		GlobalCooldown: time.Second * 4,
	}
	// load a test command
	testCommand.Load()

	type args struct {
		inMsg *messages.Incoming
	}
	tests := []struct {
		name  string
		args  args
		want  *cmd.Command
		want1 *cmd.Context
	}{
		{"empty everything", args{messages.FakeIncoming("", "", "", false, messages.Twitch)}, nil, nil},
		{"no prefix", args{messages.FakeIncoming("jtv", "Hi im a big fan!", "justinfan123", false, messages.Twitch)}, nil, nil},
		{"prefix with no command", args{messages.FakeIncoming("jtv", "|", "justinfan123", false, messages.Twitch)}, nil, nil},
		{
			"prefix with command",
			args{messages.FakeIncoming("jtv", "|testCMD", "justinfan123", false, messages.Twitch)},
			testCommand,
			&cmd.Context{
				Incoming:   *messages.FakeIncoming("jtv", "|testCMD", "justinfan123", false, messages.Twitch),
				Args:       []string{},
				Invocation: "testCMD",
			},
		},
		{
			"prefix with command and arguments",
			args{messages.FakeIncoming("jtv", "|testCMD lol xd", "justinfan123", false, messages.Twitch)},
			testCommand,
			&cmd.Context{
				Incoming:   *messages.FakeIncoming("jtv", "|testCMD lol xd", "justinfan123", false, messages.Twitch),
				Args:       []string{"lol", "xd"},
				Invocation: "testCMD",
			},
		},
		{
			"prefix with command using alias and arguments",
			args{messages.FakeIncoming("tetyys", "|AYYYYyyyyyyyLMAAAAAOOOOOO_Alien_Please AlienPls Les GOOOO", "AlienFAn", false, messages.Twitch)},
			testCommand,
			&cmd.Context{
				Incoming:   *messages.FakeIncoming("tetyys", "|AYYYYyyyyyyyLMAAAAAOOOOOO_Alien_Please AlienPls Les GOOOO", "AlienFAn", false, messages.Twitch),
				Args:       []string{"AlienPls", "Les", "GOOOO"},
				Invocation: "AYYYYyyyyyyyLMAAAAAOOOOOO_Alien_Please",
			},
		},
		{
			"prefix with command using alias with wrong capitals",
			args{messages.FakeIncoming("tetyys", "|AYYyyyyyyyyyLMAAAAAoooooo_Alien_Please AlienPls Les GOOOO", "AlienFAn", false, messages.Twitch)},
			nil, nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getCommandAndContext(tt.args.inMsg)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCommandAndContext() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getCommandAndContext() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_prepareMessage(t *testing.T) {
	type args struct {
		messageText string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"empty", args{""}, []string{}},
		{"no prefix", args{"Hi !!!"}, []string{}},
		{"prefix but no text", args{config.Public.Global.CommandPrefix}, []string{""}},
		{"prefix + cmd and 2 other arguments", args{"|testCMD lol xd"}, []string{"testCMD", "lol", "xd"}},
		{"prefix and one word", args{config.Public.Global.CommandPrefix + "yourm0M"}, []string{"yourm0M"}},
		{"prefix and multiple arguments", args{config.Public.Global.CommandPrefix + "help help FeelsDankMan how does this work?"}, []string{"help", "help", "FeelsDankMan", "how", "does", "this", "work?"}},
		{"prefix character in middle", args{"My favorite textual character(s) in the universe is '" + config.Public.Global.CommandPrefix + "'! PogChamp"}, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := prepareMessage(tt.args.messageText); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("prepareMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
