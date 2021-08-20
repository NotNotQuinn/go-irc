package command

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/NotNotQuinn/go-irc/cmd"
	"github.com/NotNotQuinn/go-irc/config"
	"github.com/NotNotQuinn/go-irc/core"
	"github.com/NotNotQuinn/go-irc/core/sender/ratelimiter"
)

var prefix string

func TestMain(m *testing.M) {
	prefix = config.Public.Global.CommandPrefix

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
		inMsg *core.Incoming
	}
	tests := []struct {
		name         string
		args         args
		wantErr      bool
		wantResponse bool
	}{
		{"nil inMsg", args{nil}, false, false},
		{"normal message", args{core.FakeIncoming("jtv", "Hi!", &core.User{
			ID:        1,
			Name:      "quinndt",
			TwitchID:  123123,
			FirstSeen: "2020-06-24 06:01:01",
		}, false, core.Twitch)}, false, false},
		{"working command", args{core.FakeIncoming("jtv", prefix+"working lol", &core.User{
			ID:        1,
			Name:      "quinndt",
			TwitchID:  123123,
			FirstSeen: "2020-06-24 06:01:01",
		}, false, core.Twitch)}, false, true},
		{"working command with alias", args{core.FakeIncoming("jtv", prefix+"Work lol", &core.User{
			ID:        1,
			Name:      "quinndt",
			TwitchID:  123123,
			FirstSeen: "2020-06-24 06:01:01",
		}, false, core.Twitch)}, false, true},
		{"erroring command with response", args{core.FakeIncoming("jtv", prefix+"Error lol xd", &core.User{
			ID:        1,
			Name:      "quinndt",
			TwitchID:  123123,
			FirstSeen: "2020-06-24 06:01:01",
		}, false, core.Twitch)}, true, true},
		{"erroring command without response", args{core.FakeIncoming("jtv", prefix+"Error xd", &core.User{
			ID:        1,
			Name:      "quinndt",
			TwitchID:  123123,
			FirstSeen: "2020-06-24 06:01:01",
		}, false, core.Twitch)}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := HandleMessage(tt.args.inMsg); (err != nil) != tt.wantErr {
				t.Errorf("HandleMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
			var res *core.Outgoing
			time.Sleep(time.Second / 100)
			select {
			case res = <-core.MessagesOUT:
			default:
			}
			if (res != nil) != tt.wantResponse {
				t.Errorf("HandleMessage(); <-core.MessagesOUT %v, wantResponse %v", res, tt.wantResponse)
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
		inMsg *core.Incoming
	}
	tests := []struct {
		name string
		args args
		want *cmd.Context
	}{
		{"nil msg", args{nil}, nil},
		{
			"empty everything", args{core.FakeIncoming("", "", &core.User{ID: 0, Name: "", TwitchID: 0, FirstSeen: ""}, false, core.Twitch)},
			&cmd.Context{
				Incoming:   core.FakeIncoming("", "", &core.User{ID: 0, Name: "", TwitchID: 0, FirstSeen: ""}, false, core.Twitch),
				Args:       []string{},
				Invocation: "",
				User:       &core.User{ID: 0, Name: "", TwitchID: 0, FirstSeen: ""},
			},
		},
		{
			"no prefix", args{core.FakeIncoming("jtv", "Hi im a big fan!", &core.User{
				ID:        10,
				Name:      "justinfan123",
				TwitchID:  10101010101123,
				FirstSeen: time.Date(2007, 0, 0, 0, 0, 0, 0, time.UTC).Format("2006-01-02 03:04:05"),
			}, false, core.Twitch)},
			&cmd.Context{
				Incoming: core.FakeIncoming("jtv", "Hi im a big fan!", &core.User{
					ID:        10,
					Name:      "justinfan123",
					TwitchID:  10101010101123,
					FirstSeen: time.Date(2007, 0, 0, 0, 0, 0, 0, time.UTC).Format("2006-01-02 03:04:05"),
				}, false, core.Twitch),
				Args:       []string{"im", "a", "big", "fan!"},
				Invocation: "Hi",
				User: &core.User{
					ID:        10,
					Name:      "justinfan123",
					TwitchID:  10101010101123,
					FirstSeen: time.Date(2007, 0, 0, 0, 0, 0, 0, time.UTC).Format("2006-01-02 03:04:05"),
				},
			},
		},
		{
			"prefix with no command",
			args{core.FakeIncoming("jtv", prefix, &core.User{
				ID:        10,
				Name:      "justinfan123",
				TwitchID:  10101010101123,
				FirstSeen: time.Date(2007, 0, 0, 0, 0, 0, 0, time.UTC).Format("2006-01-02 03:04:05"),
			}, false, core.Twitch)},
			&cmd.Context{
				Incoming: core.FakeIncoming("jtv", prefix, &core.User{
					ID:        10,
					Name:      "justinfan123",
					TwitchID:  10101010101123,
					FirstSeen: time.Date(2007, 0, 0, 0, 0, 0, 0, time.UTC).Format("2006-01-02 03:04:05"),
				}, false, core.Twitch),
				Args:       []string{},
				Invocation: "",
				Command:    nil,
				User: &core.User{
					ID:        10,
					Name:      "justinfan123",
					TwitchID:  10101010101123,
					FirstSeen: time.Date(2007, 0, 0, 0, 0, 0, 0, time.UTC).Format("2006-01-02 03:04:05"),
				},
			},
		},
		{
			"prefix with command",
			args{core.FakeIncoming("jtv", prefix+"testCMD", &core.User{
				ID:        10,
				Name:      "justinfan123",
				TwitchID:  10101010101123,
				FirstSeen: time.Date(2007, 0, 0, 0, 0, 0, 0, time.UTC).Format("2006-01-02 03:04:05"),
			}, false, core.Twitch)},
			&cmd.Context{
				Incoming: core.FakeIncoming("jtv", prefix+"testCMD", &core.User{
					ID:        10,
					Name:      "justinfan123",
					TwitchID:  10101010101123,
					FirstSeen: time.Date(2007, 0, 0, 0, 0, 0, 0, time.UTC).Format("2006-01-02 03:04:05"),
				}, false, core.Twitch),
				Args:       []string{},
				Invocation: "testCMD",
				Command:    testCommand,
				User: &core.User{
					ID:        10,
					Name:      "justinfan123",
					TwitchID:  10101010101123,
					FirstSeen: time.Date(2007, 0, 0, 0, 0, 0, 0, time.UTC).Format("2006-01-02 03:04:05"),
				},
			},
		},
		{
			"prefix with command and arguments",
			args{core.FakeIncoming("jtv", prefix+"testCMD lol xd", &core.User{
				ID:        10,
				Name:      "justinfan123",
				TwitchID:  10101010101123,
				FirstSeen: time.Date(2007, 0, 0, 0, 0, 0, 0, time.UTC).Format("2006-01-02 03:04:05"),
			}, false, core.Twitch)},
			&cmd.Context{
				Incoming: core.FakeIncoming("jtv", prefix+"testCMD lol xd", &core.User{
					ID:        10,
					Name:      "justinfan123",
					TwitchID:  10101010101123,
					FirstSeen: time.Date(2007, 0, 0, 0, 0, 0, 0, time.UTC).Format("2006-01-02 03:04:05"),
				}, false, core.Twitch),
				Args:       []string{"lol", "xd"},
				Invocation: "testCMD",
				Command:    testCommand,
				User: &core.User{
					ID:        10,
					Name:      "justinfan123",
					TwitchID:  10101010101123,
					FirstSeen: time.Date(2007, 0, 0, 0, 0, 0, 0, time.UTC).Format("2006-01-02 03:04:05"),
				},
			},
		},
		{
			"prefix with command using alias and arguments",
			args{core.FakeIncoming("tetyys", prefix+"AYYYYyyyyyyyLMAAAAAOOOOOO_Alien_Please AlienPls Les GOOOO", &core.User{
				ID:        1000,
				Name:      "AlienFAn",
				TwitchID:  777777,
				FirstSeen: time.Date(2017, 10, 10, 21, 1, 1, 0, time.UTC).Format("2006-01-02 03:04:05"),
			}, false, core.Twitch)},
			&cmd.Context{
				Incoming: core.FakeIncoming("tetyys", prefix+"AYYYYyyyyyyyLMAAAAAOOOOOO_Alien_Please AlienPls Les GOOOO", &core.User{
					ID:        1000,
					Name:      "AlienFAn",
					TwitchID:  777777,
					FirstSeen: time.Date(2017, 10, 10, 21, 1, 1, 0, time.UTC).Format("2006-01-02 03:04:05"),
				}, false, core.Twitch),
				Args:       []string{"AlienPls", "Les", "GOOOO"},
				Invocation: "AYYYYyyyyyyyLMAAAAAOOOOOO_Alien_Please",
				Command:    testCommand,
				User: &core.User{
					ID:        1000,
					Name:      "AlienFAn",
					TwitchID:  777777,
					FirstSeen: time.Date(2017, 10, 10, 21, 1, 1, 0, time.UTC).Format("2006-01-02 03:04:05"),
				},
			},
		},
		{
			"prefix with command using alias with wrong capitals",
			args{core.FakeIncoming("tetyys", prefix+"AYYyyyyyyyyyLMAAAAAoooooo_Alien_Please AlienPls Les GOOOO", &core.User{
				ID:        1000,
				Name:      "AlienFAn",
				TwitchID:  777777,
				FirstSeen: time.Date(2017, 10, 10, 21, 1, 1, 0, time.UTC).Format("2006-01-02 03:04:05"),
			}, false, core.Twitch)},
			&cmd.Context{
				Incoming: core.FakeIncoming("tetyys", prefix+"AYYyyyyyyyyyLMAAAAAoooooo_Alien_Please AlienPls Les GOOOO", &core.User{
					ID:        1000,
					Name:      "AlienFAn",
					TwitchID:  777777,
					FirstSeen: time.Date(2017, 10, 10, 21, 1, 1, 0, time.UTC).Format("2006-01-02 03:04:05"),
				}, false, core.Twitch),
				Args:       []string{"AlienPls", "Les", "GOOOO"},
				Invocation: "AYYyyyyyyyyyLMAAAAAoooooo_Alien_Please",
				Command:    nil,
				User: &core.User{
					ID:        1000,
					Name:      "AlienFAn",
					TwitchID:  777777,
					FirstSeen: time.Date(2017, 10, 10, 21, 1, 1, 0, time.UTC).Format("2006-01-02 03:04:05"),
				},
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
		{"prefix but no text", args{prefix}, []string{""}, true},
		{"prefix + cmd and 2 other arguments", args{prefix + "testCMD lol xd"}, []string{"testCMD", "lol", "xd"}, true},
		{"prefix and one word", args{prefix + "yourm0M"}, []string{"yourm0M"}, true},
		{"prefix and multiple arguments", args{prefix + "help help FeelsDankMan how does this work?"}, []string{"help", "help", "FeelsDankMan", "how", "does", "this", "work?"}, true},
		{"prefix character in middle", args{"My favorite textual character(s) in the universe is '" + prefix + "'! PogChamp"}, []string{
			"My", "favorite", "textual", "character(s)", "in", "the", "universe", "is", "'" + prefix + "'!", "PogChamp",
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
