package core

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/NotNotQuinn/go-irc/config"
	"github.com/NotNotQuinn/go-irc/data"
)

var (
	// Maps user (internal) ID to user pointer.
	userIDMap = make(map[uint64]*User)
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

// Represents ability to perform actions as a specific user
type UserPermissions struct {
	// The user who's permissions are represented
	User *User
	// Whether the user is admin
	Admin bool
}

// Get the permissions of the user
func (u *User) Perms() *UserPermissions {
	perms := DefaultPermissions()
	if config.Public.Users.Admins.Inclues(u.Name) {
		perms.Admin = true
	}
	return perms
}

// Get a user.
func GetUser(Name string, ID uint64) (*User, error) {
	if data.CoreDB == nil {
		return nil, errors.New("database not availible")
	}
	var u *User
	var cond string
	var args []interface{}
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
func AlwaysGetUser(Name string, TwitchID uint64) *User {
	u, err := GetOrCreateUser(Name, TwitchID)
	if err != nil {
		Errors <- fmt.Errorf("error getting user for message: %w", err)
		u = &User{
			TwitchID:  TwitchID,
			ID:        0,
			Name:      Name,
			FirstSeen: "",
		}
	}
	return u
}

// Gets a user, or creates it.
func GetOrCreateUser(Name string, TwitchID uint64) (*User, error) {
	u, err := GetUser(Name, 0)
	if err != nil {
		u, err = CreateNewUser(Name, TwitchID)
		if err != nil {
			return nil, err
		}
	}
	if u.TwitchID != TwitchID {
		return nil, errors.New("username does not match twitch id")
	}

	return u, nil
}

// Creates a new user, if the user already exists return nil.
func CreateNewUser(Name string, Twitch_ID uint64) (*User, error) {
	if data.CoreDB == nil {
		return nil, errors.New("database not availible")
	}
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
		now.Format("2006-01-02 03:04:05"),
	}, nil
}

// Gets the permissions of a user with no permissions
func DefaultPermissions() *UserPermissions {
	return &UserPermissions{
		Admin: false,
		User:  &User{ID: 0, Name: "", TwitchID: 0, FirstSeen: ""},
	}
}
