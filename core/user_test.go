package core

import (
	"reflect"
	"testing"
)

func TestGetUser(t *testing.T) {
	type args struct {
		Name string
		ID   uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *User
		wantErr bool
	}{
		{
			name: "Only provide name",
			args: args{Name: "quinndt"},
			want: &User{ID: 1, Name: "quinndt", TwitchID: 440674731, FirstSeen: "2020-06-24 06:01:01"},
		},
		{
			name: "Only provide id",
			args: args{ID: 1},
			want: &User{ID: 1, Name: "quinndt", TwitchID: 440674731, FirstSeen: "2020-06-24 06:01:01"},
		},
		{
			name: "Provide name & id",
			args: args{Name: "quinndt", ID: 1},
			want: &User{ID: 1, Name: "quinndt", TwitchID: 440674731, FirstSeen: "2020-06-24 06:01:01"},
		},
		{
			name:    "Provide mismatched name & id",
			args:    args{Name: "turtoise", ID: 1},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Name that does not exist.",
			args:    args{Name: "_not_a_user"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Id that does not exist.",
			args:    args{ID: 9387597663878},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Name and ID that does not exist.",
			args:    args{Name: "_not_a_user_2", ID: 346536888863},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUser(tt.args.Name, tt.args.ID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlwaysGetUser(t *testing.T) {
	type args struct {
		Name     string
		TwitchID uint64
	}
	tests := []struct {
		name string
		args args
		want *User
	}{
		// TODO: Add test cases, if we keep this
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AlwaysGetUser(tt.args.Name, tt.args.TwitchID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AlwaysGetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetOrCreateUser(t *testing.T) {
	type args struct {
		Name     string
		TwitchID uint64
	}
	tests := []struct {
		name               string
		args               args
		want               *User
		skipFirstSeenAndID bool
		wantErr            bool
	}{
		{
			name: "Get user",
			args: args{Name: "quinndt", TwitchID: 440674731},
			want: &User{ID: 1, Name: "quinndt", TwitchID: 440674731, FirstSeen: "2020-06-24 06:01:01"},
		},
		{
			name:               "Create new user.",
			args:               args{Name: "_create_GetOrCreateUser_1", TwitchID: 999389473},
			want:               &User{ID: 0, Name: "_create_GetOrCreateUser_1", TwitchID: 999389473, FirstSeen: ""},
			skipFirstSeenAndID: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetOrCreateUser(tt.args.Name, tt.args.TwitchID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOrCreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				if tt.skipFirstSeenAndID && tt.want.Name == got.Name && tt.want.TwitchID == got.TwitchID {
					// This is ok, because we dont check the firstseen and local id
					// (we created it just now, so how would we know?!)
					// Other values are equal though!
					return
				}
				t.Errorf("GetOrCreateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateNewUser(t *testing.T) {
	type args struct {
		Name      string
		Twitch_ID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *User
		wantErr bool
	}{
		{
			name:    "user exists",
			args:    args{Name: "quinndt", Twitch_ID: 440674731},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "duplicate name",
			args:    args{Name: "turtoise", Twitch_ID: 397447},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "duplicate id",
			args:    args{Name: "_urmom", Twitch_ID: 440674731},
			want:    nil,
			wantErr: true,
		},
		{
			name: "create new user",
			args: args{Name: "_create_CreateNewUser_1", Twitch_ID: 3128405658},
			want: &User{Name: "_create_CreateNewUser_1", TwitchID: 3128405658},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateNewUser(tt.args.Name, tt.args.Twitch_ID)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateNewUser() error = %v, wantErr %v", err, tt.wantErr)
				t.Errorf("CreateNewUser() = %v, want %v", got, tt.want)
				return
			}
			// check because could be nil
			if reflect.DeepEqual(got, tt.want) {
				return
			}
			// Skip checking ID and FirstSeen because those are variable in this case.
			if got.TwitchID != tt.want.TwitchID && got.Name != tt.want.Name {
				t.Errorf("CreateNewUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
