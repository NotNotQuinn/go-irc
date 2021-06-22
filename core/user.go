package core

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/NotNotQuinn/go-irc/config"
	"github.com/NotNotQuinn/go-irc/data"
	"github.com/gempir/go-twitch-irc/v2"
)

var (
	// Maps user (internal) ID to user pointer.
	userIDMap = make(map[uint]*User)
	// Maps user name to user pointer.
	userNameMap = make(map[string]*User)
)

// A user
type User struct {
	ID        uint
	Name      string
	TwitchID  uint64
	FirstSeen string
}

func FakeUser(s string) User {
	return User{0, s, 0, ""}
}

// Represents ability to perform actions as a specific user
type UserPermissions struct {
	// The user who's permissions are represented
	User User
	// Whether the user is admin
	Admin bool
}

// Get the permissions of the user
func (u User) GetPermissions() *UserPermissions {
	perms := DefaultPermissions()
	if config.Public.Users.Admins.Inclues(u.Name) {
		perms.Admin = true
	}
	return perms
}

// Get a user.
//
// Allowed types: string, and int
func GetUser(Name string, ID uint) (*User, error) {
	var u *User
	var cond string
	var args []interface{} = make([]interface{}, 2)
	if Name != "" && ID != 0 {
		u = userIDMap[ID]
		u2 := userNameMap[Name]
		if !reflect.DeepEqual(u, u2) {
			return nil, fmt.Errorf("user not found: user Name and ID point to seperate users")
		}
		cond = "Name = ? AND ID = ?"
		args = []interface{}{Name, ID}
	} else if Name == "" && ID != 0 {
		u = userIDMap[ID]
		cond = "ID = ?"
		args = []interface{}{ID}
	} else if Name != "" && ID == 0 {
		u = userNameMap[Name]
		cond = "Name = ?"
		args = []interface{}{Name}
	} else {
		return nil, fmt.Errorf("not enough information to get user (name and id provided as 0 values)")
	}
	// Have both
	if u == nil {
		row := data.CoreDB.QueryRow("SELECT ID, Name, Twitch_ID, First_Seen FROM wb_core.user WHERE ("+cond+");", args...)
		user := &User{}
		err := row.Scan(&user.ID, &user.Name, &user.TwitchID, &user.FirstSeen)
		if err != nil {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		u = user
	}
	return u, nil
}

// Will almost never fail to get a user
//
// Falls back to basic data from twitch user, if there was a problem getting or creating the user.
func AlwaysGetUser(tu twitch.User) *User {
	u, err := getUser(tu)
	if err != nil {
		Errors <- fmt.Errorf("error getting user for message: %w", err)
		TID, err := strconv.ParseUint(tu.ID, 10, 1)
		if err != nil {
			// Should never happen
			Errors <- err
			return nil
		}
		u = &User{
			TwitchID:  TID,
			ID:        0,
			Name:      tu.Name,
			FirstSeen: "",
		}
	}
	return u
}

// Gets a user
func getUser(tu twitch.User) (*User, error) {
	u, err := GetUser(tu.Name, 0)
	if err != nil {
		return nil, err
	}
	TID, err := strconv.ParseUint(tu.ID, 10, 1)
	if err != nil {
		return nil, err
	}
	if u.TwitchID != TID {
		return nil, errors.New("username does not match twitch id")
	}

	if u == nil {
		u, err = CreateNewUser(tu.Name, TID)
		if err != nil {
			return nil, err
		}
	}

	return u, nil
}

// Creates a new user, if the user already exists return nil.
func CreateNewUser(Name string, Twitch_ID uint64) (*User, error) {
	if u, _ := GetUser(Name, 0); u != nil && u.TwitchID == Twitch_ID {
		return nil, fmt.Errorf("user not created: already exists with ID %d", u.ID)
	}
	now := time.Now()
	res, err := data.CoreDB.Exec("INSERT INTO `wb_core`.`user` (`Name`, `Twitch_ID`, `First_Seen`) VALUES (?, ?, ?)", Name, Twitch_ID, now)
	if err != nil {
		return nil, fmt.Errorf("user not created: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("user ID not fetched: %w", err)
	}
	return &User{
		uint(id),
		Name,
		uint64(Twitch_ID),
		now.String(),
	}, nil
}

// Gets the permissions of a user with no permissions
func DefaultPermissions() *UserPermissions {
	return &UserPermissions{
		Admin: false,
		User:  FakeUser(""),
	}
}
