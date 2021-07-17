package core

import (
	"reflect"
	"testing"

	"github.com/gempir/go-twitch-irc/v2"
)

func TestNewIncoming(t *testing.T) {
	type args struct {
		msg interface{ GetType() twitch.MessageType }
	}
	// Must be alocated beforehand because the raw field accepts a pointer only.
	// I did it this way because I dont see any other easy way to generate fake twitch messages from the library.
	msg_michaelreeves_quinndt := twitch.ParseMessage("@badge-info=subscriber/8;badges=moderator/1,subscriber/6,bits/1;color=#B1FCDF;display-name=QuinnDT;emotes=;flags=;id=01619e22-5b6b-47df-aa33-d974a4980faf;mod=1;rm-deleted=1;rm-received-ts=1623998086287;room-id=469790580;subscriber=1;tmi-sent-ts=1623998115148;turbo=0;user-id=440674731;user-type=mod :quinndt!quinndt@quinndt.tmi.twitch.tv PRIVMSG #michaelreeves Pog")
	msg_supinic_quinndt := twitch.ParseMessage("@badge-info=;badges=;color=#B1FCDF;display-name=QuinnDT;emotes=;flags=;id=342af4f3-ebb8-46ef-9cdf-b71caf05780a;mod=0;rm-received-ts=1624050851701;room-id=31400525;subscriber=0;tmi-sent-ts=1624050851544;turbo=0;user-id=440674731;user-type= :quinndt!quinndt@quinndt.tmi.twitch.tv PRIVMSG #supinic :APU test 1 2 3 2 1 tset upA")
	whisper_quinndt := twitch.ParseMessage("@badges=;color=#B1FCDF;display-name=QuinnDT;emotes=;message-id=1038;thread-id=564777265_440674731;turbo=0;user-id=440674731;user-type= :quinndt!quinndt@quinndt.tmi.twitch.tv WHISPER wanductbot :Hi :)")
	tests := []struct {
		name string
		args args
		want *Incoming
	}{
		{"nil inmsg", args{nil}, nil},
		{"Moderator and sub message", args{msg_michaelreeves_quinndt}, &Incoming{
			Platform: Twitch,
			Channel:  "michaelreeves",
			Message:  "Pog",
			User: &User{
				1,
				"quinndt",
				440674731,
				"2020-06-24 06:01:01",
			},
			Raw: &msg_michaelreeves_quinndt,
			DMs: false,
		}},
		{"non-sub pleb message", args{msg_supinic_quinndt}, &Incoming{
			Platform: Twitch,
			Channel:  "supinic",
			Message:  "APU test 1 2 3 2 1 tset upA",
			User: &User{
				1,
				"quinndt",
				440674731,
				"2020-06-24 06:01:01",
			},
			Raw: &msg_supinic_quinndt,
			DMs: false,
		}},
		{"whisper from quinndt", args{whisper_quinndt}, &Incoming{
			Platform: Twitch,
			Channel:  "",
			Message:  "Hi :)",
			User: &User{
				1,
				"quinndt",
				440674731,
				"2020-06-24 06:01:01",
			},
			Raw: &whisper_quinndt,
			DMs: true,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewIncoming(tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIncoming() = %v, want %v", got, tt.want)
				t.Errorf("NewIncoming().User = %v, want %v", got.User, tt.want.User)
			}
		})
	}
}

func TestNewOutgoing(t *testing.T) {
	type args struct {
		inMsg    *Incoming
		responce string
	}
	tests := []struct {
		name string
		args args
		want *Outgoing
	}{
		{"nil inMsg", args{nil, "This should not matter..."}, nil},
		{"message in channel", args{&Incoming{
			Platform: Twitch,
			Channel:  "quinndt",
			Message:  "yo",
			User: &User{
				1,
				"quinndt",
				440674731,
				"2020-06-24 06:01:01",
			},
			Raw: nil,
			DMs: false,
		}, "WADUP!!"}, &Outgoing{
			Platform: Twitch,
			Message:  "WADUP!!",
			Channel:  "quinndt",
			User: &User{
				1,
				"quinndt",
				440674731,
				"2020-06-24 06:01:01",
			},
			DM:       false,
			NoFilter: false,
		}},
		{"whisper message", args{&Incoming{
			Platform: Twitch,
			Channel:  "",
			Message:  "Hi!!",
			User:     &User{2, "turtoise", 80805824, "2020-06-24 06:01:01"},
			Raw:      nil,
			DMs:      true,
		}, "yo"}, &Outgoing{
			Platform: Twitch,
			Message:  "yo",
			Channel:  "",
			User:     &User{2, "turtoise", 80805824, "2020-06-24 06:01:01"},
			DM:       true,
			NoFilter: false,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOutgoing(tt.args.inMsg, tt.args.responce); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOutgoing() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFakeOutgoing(t *testing.T) {
	type args struct {
		channel  string
		message  string
		platform PlatformType
	}
	tests := []struct {
		name string
		args args
		want *Outgoing
	}{
		{"nil or empty everything", args{"", "", Unknown}, &Outgoing{
			Platform: Unknown,
			Message:  "",
			Channel:  "",
			User:     &User{ID: 0, Name: "", TwitchID: 0, FirstSeen: ""},
			DM:       false,
			NoFilter: false,
		}},
		{"channel and user on twitch platform", args{"quinndt", "yoo!", Twitch}, &Outgoing{
			Platform: Twitch,
			Message:  "yoo!",
			Channel:  "quinndt",
			User:     &User{ID: 0, Name: "", TwitchID: 0, FirstSeen: ""},
			DM:       false,
			NoFilter: false,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FakeOutgoing(tt.args.channel, tt.args.message, tt.args.platform); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FakeOutgoing() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFakeIncoming(t *testing.T) {
	type args struct {
		channel  string
		message  string
		user     *User
		DMs      bool
		platform PlatformType
	}
	tests := []struct {
		name string
		args args
		want *Incoming
	}{
		{"least information possible", args{"", "", &User{ID: 0, Name: "", TwitchID: 0, FirstSeen: ""}, false, Unknown}, &Incoming{
			Platform: Unknown,
			Channel:  "",
			Message:  "",
			User:     &User{ID: 0, Name: "", TwitchID: 0, FirstSeen: ""},
			Raw:      nil,
			DMs:      false,
		}},
		{"normal message", args{"quinndt", "squadR turdoise", &User{2, "turtoise", 80805824, "2020-06-24 06:01:01"}, false, Twitch}, &Incoming{
			Platform: Twitch,
			Channel:  "quinndt",
			Message:  "squadR turdoise",
			User:     &User{2, "turtoise", 80805824, "2020-06-24 06:01:01"},
			Raw:      nil,
			DMs:      false,
		}},
		{"whispers", args{"", "pog it didnt crash", &User{
			1,
			"quinndt",
			123123,
			"2020-06-24 06:01:01",
		}, true, Twitch}, &Incoming{
			Platform: Twitch,
			Channel:  "",
			Message:  "pog it didnt crash",
			User: &User{
				1,
				"quinndt",
				123123,
				"2020-06-24 06:01:01",
			},
			Raw: nil,
			DMs: true,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FakeIncoming(tt.args.channel, tt.args.message, tt.args.user, tt.args.DMs, tt.args.platform); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FakeIncoming() = %v, want %v", got, tt.want)
			}
		})
	}
}
