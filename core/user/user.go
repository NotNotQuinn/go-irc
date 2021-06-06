package wbUser

import "github.com/NotNotQuinn/go-irc/config"

// A user
type User string

type IUser interface {
	GetPermissions() *Permissions
	Name() string
}

// Represents ability to perform actions as a specific user
type Permissions struct {
	// The user who's permissions are represented
	User User
	// Whether the user is admin
	Admin bool
}

// Get the permissions of the user
func (u User) GetPermissions() *Permissions {
	perms := DefaultPermissions()
	if config.Public.Users.Admins.Inclues(string(u)) {
		perms.Admin = true
	}
	return perms
}

// The name of the user
func (u User) Name() string {
	return string(u)
}

// Gets the permissions of a user with no permissions
func DefaultPermissions() *Permissions {
	return &Permissions{
		Admin: false,
		User:  "",
	}
}
