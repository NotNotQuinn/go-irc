package command

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/NotNotQuinn/go-irc/channels"
	"github.com/NotNotQuinn/go-irc/cmd"
	"github.com/NotNotQuinn/go-irc/config"
	"github.com/NotNotQuinn/go-irc/core/command/messages"
	"github.com/NotNotQuinn/go-irc/core/sender/ratelimiter"
	wbUser "github.com/NotNotQuinn/go-irc/core/user"
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
	var erroringCommand = &cmd.Command{
		Name:    "erroring",
		Aliases: []string{"err", "Error"},
		Execution: func(c *cmd.Context) (*cmd.Return, error) {
			if c.Args[0] == "lol" {
				return &cmd.Return{
					Success: true,
					Reply:   "xd",
				}, fmt.Errorf("generic error")
			}
			return nil, fmt.Errorf("generic error")
		},
	}
	var workingCommand = &cmd.Command{
		Name:    "working",
		Aliases: []string{"Work"},
		Execution: func(c *cmd.Context) (*cmd.Return, error) {
			return &cmd.Return{
				Success: true,
				Reply:   "Hi!",
			}, nil
		},
	}
	erroringCommand.Load()
	workingCommand.Load()
	type args struct {
		inMsg *messages.Incoming
	}
	tests := []struct {
		name         string
		args         args
		wantErr      bool
		wantResponse bool
	}{
		{"nil inMsg", args{nil}, false, false},
		{"normal message", args{messages.FakeIncoming("jtv", "Hi!", wbUser.FakeUser("quinndt"), false, messages.Twitch)}, false, false},
		{"working command", args{messages.FakeIncoming("jtv", "|working lol", wbUser.FakeUser("quinndt"), false, messages.Twitch)}, false, true},
		{"working command with alias", args{messages.FakeIncoming("jtv", "|Work lol", wbUser.FakeUser("quinndt"), false, messages.Twitch)}, false, true},
		{"erroring command with response", args{messages.FakeIncoming("jtv", "|Error lol xd", wbUser.FakeUser("quinndt"), false, messages.Twitch)}, true, true},
		{"erroring command without response", args{messages.FakeIncoming("jtv", "|Error xd", wbUser.FakeUser("quinndt"), false, messages.Twitch)}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := HandleMessage(tt.args.inMsg); (err != nil) != tt.wantErr {
				t.Errorf("HandleMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
			var res *messages.Outgoing
			time.Sleep(time.Second / 100)
			select {
			case res = <-channels.MessagesOUT:
			default:
			}
			if (res != nil) != tt.wantResponse {
				t.Errorf("HandleMessage(); <-channels.MessagesOUT = %v, wantResponse %v", res, tt.wantResponse)
			}
		})
	}
}

func TestGetContext(t *testing.T) {
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
		name string
		args args
		want *cmd.Context
	}{
		{"nil msg", args{nil}, nil},
		{
			"empty everything", args{messages.FakeIncoming("", "", wbUser.FakeUser(""), false, messages.Twitch)},
			&cmd.Context{
				Incoming:   messages.FakeIncoming("", "", wbUser.FakeUser(""), false, messages.Twitch),
				Args:       []string{},
				Invocation: "",
			},
		},
		{
			"no prefix", args{messages.FakeIncoming("jtv", "Hi im a big fan!", wbUser.FakeUser("justinfan123"), false, messages.Twitch)},
			&cmd.Context{
				Incoming:   messages.FakeIncoming("jtv", "Hi im a big fan!", wbUser.FakeUser("justinfan123"), false, messages.Twitch),
				Args:       []string{"im", "a", "big", "fan!"},
				Invocation: "Hi",
			},
		},
		{
			"prefix with no command",
			args{messages.FakeIncoming("jtv", "|", wbUser.FakeUser("justinfan123"), false, messages.Twitch)},
			&cmd.Context{
				Incoming:   messages.FakeIncoming("jtv", "|", wbUser.FakeUser("justinfan123"), false, messages.Twitch),
				Args:       []string{},
				Invocation: "",
				Command:    nil,
			},
		},
		{
			"prefix with command",
			args{messages.FakeIncoming("jtv", "|testCMD", wbUser.FakeUser("justinfan123"), false, messages.Twitch)},
			&cmd.Context{
				Incoming:   messages.FakeIncoming("jtv", "|testCMD", wbUser.FakeUser("justinfan123"), false, messages.Twitch),
				Args:       []string{},
				Invocation: "testCMD",
				Command:    testCommand,
			},
		},
		{
			"prefix with command and arguments",
			args{messages.FakeIncoming("jtv", "|testCMD lol xd", wbUser.FakeUser("justinfan123"), false, messages.Twitch)},
			&cmd.Context{
				Incoming:   messages.FakeIncoming("jtv", "|testCMD lol xd", wbUser.FakeUser("justinfan123"), false, messages.Twitch),
				Args:       []string{"lol", "xd"},
				Invocation: "testCMD",
				Command:    testCommand,
			},
		},
		{
			"prefix with command using alias and arguments",
			args{messages.FakeIncoming("tetyys", "|AYYYYyyyyyyyLMAAAAAOOOOOO_Alien_Please AlienPls Les GOOOO", wbUser.FakeUser("AlienFAn"), false, messages.Twitch)},
			&cmd.Context{
				Incoming:   messages.FakeIncoming("tetyys", "|AYYYYyyyyyyyLMAAAAAOOOOOO_Alien_Please AlienPls Les GOOOO", wbUser.FakeUser("AlienFAn"), false, messages.Twitch),
				Args:       []string{"AlienPls", "Les", "GOOOO"},
				Invocation: "AYYYYyyyyyyyLMAAAAAOOOOOO_Alien_Please",
				Command:    testCommand,
			},
		},
		{
			"prefix with command using alias with wrong capitals",
			args{messages.FakeIncoming("tetyys", "|AYYyyyyyyyyyLMAAAAAoooooo_Alien_Please AlienPls Les GOOOO", wbUser.FakeUser("AlienFAn"), false, messages.Twitch)},
			&cmd.Context{
				Incoming:   messages.FakeIncoming("tetyys", "|AYYyyyyyyyyyLMAAAAAoooooo_Alien_Please AlienPls Les GOOOO", wbUser.FakeUser("AlienFAn"), false, messages.Twitch),
				Args:       []string{"AlienPls", "Les", "GOOOO"},
				Invocation: "AYYyyyyyyyyyLMAAAAAoooooo_Alien_Please",
				Command:    nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := GetContext(tt.args.inMsg)
			if !reflect.DeepEqual(ctx, tt.want) {
				t.Errorf("getContext() = %v, want %v", ctx, tt.want)
			}
		})
	}
}

func Test_prepareMessage(t *testing.T) {
	type args struct {
		messageText string
	}
	tests := []struct {
		name       string
		args       args
		want_args  []string
		want_isCMD bool
	}{
		{"empty", args{""}, []string{""}, false},
		{"no prefix", args{"Hi !!!"}, []string{"Hi", "!!!"}, false},
		{"prefix but no text", args{config.Public.Global.CommandPrefix}, []string{""}, true},
		{"prefix + cmd and 2 other arguments", args{"|testCMD lol xd"}, []string{"testCMD", "lol", "xd"}, true},
		{"prefix and one word", args{config.Public.Global.CommandPrefix + "yourm0M"}, []string{"yourm0M"}, true},
		{"prefix and multiple arguments", args{config.Public.Global.CommandPrefix + "help help FeelsDankMan how does this work?"}, []string{"help", "help", "FeelsDankMan", "how", "does", "this", "work?"}, true},
		{"prefix character in middle", args{"My favorite textual character(s) in the universe is '" + config.Public.Global.CommandPrefix + "'! PogChamp"}, []string{
			"My", "favorite", "textual", "character(s)", "in", "the", "universe", "is", "'" + config.Public.Global.CommandPrefix + "'!", "PogChamp",
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isCMD, args := prepareMessage(tt.args.messageText)
			if !reflect.DeepEqual(args, tt.want_args) {
				t.Errorf("prepareMessage() args = %v, want %v", args, tt.want_args)
			}
			if !reflect.DeepEqual(isCMD, tt.want_isCMD) {
				t.Errorf("prepareMessage() isCMD = %v, want %v", isCMD, tt.want_isCMD)
			}
		})
	}
}
