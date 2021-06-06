package config

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
)

// Config file to read from
var confFile = "./config/public_conf.json"

// Public config data
var Public *PublicConfig

// Private config data
var Private *PrivateConfig

// Public twitch config data
type PublicTwitchConfig struct {
	Channels stringList `json:"channels"`
}

// Public config data that affects globally
type PublicGlobalConfig struct {
	CommandPrefix string `json:"commandPrefix"`
	UserAgent     string `json:"user_agent"`
}

// Public data about users.
type PublicUsersConfig struct {
	Admins stringList `json:"admins"`
}

// A convinent way to interact with string slices
type stringList []string

// Check if the list includes the item provided
func (l stringList) Inclues(query string) bool {
	for _, item := range l {
		if item == query {
			return true
		}
	}
	return false
}

// All public config data
type PublicConfig struct {
	Twitch PublicTwitchConfig `json:"twitch"`
	Global PublicGlobalConfig `json:"global"`
	Users  PublicUsersConfig  `json:"users"`
}

// Load the configs and assign them to the variables
func Init() error {
	var err error
	Private, err = getPrivate()
	if err != nil {
		return err
	}
	Public, err = getPublic()
	if err != nil {
		return err
	}
	return nil
}

// Load and return the public config
func getPublic() (*PublicConfig, error) {
	bytes, err := ioutil.ReadFile(confFile)
	if err != nil {
		return nil, err
	}
	var config PublicConfig
	json.Unmarshal(bytes, &config)
	return &config, nil
}

// Reload the config from file
func (conf *PublicConfig) Reload() error {
	bytes, err := ioutil.ReadFile(confFile)
	if err != nil {
		return err
	}
	var config PublicConfig
	json.Unmarshal(bytes, &config)
	return nil
}

// Save the config to the file
func (conf *PublicConfig) Save() (success bool, err error) {
	bytes, err := json.MarshalIndent(conf, "", "\t")
	if err != nil {
		return false, err
	}
	file, err := os.Create(confFile)
	if err != nil {
		return false, err
	}
	numBytes, err := file.Write(bytes)
	if err != nil {
		file.Close()
		return false, err
	}
	fmt.Printf("Conf saved successfully (%d bytes)\n", numBytes)
	err = file.Close()
	if err != nil {
		return true, err
	}
	err = ioutil.WriteFile(confFile, bytes, fs.ModeType)
	return true, err
}

// All private config data
type PrivateConfig struct {
	// Even though the username is not 'private' its very much related so I will keep it here

	// Username of account
	Username string `json:"username"`
	// Oauth token of account
	Oauth string `json:"oauth"`
}

// Load and return private config
func getPrivate() (conf *PrivateConfig, err error) {
	bytes, err := ioutil.ReadFile("./config/private_conf.json")
	if err != nil {
		return nil, err
	}

	var config PrivateConfig
	err = json.Unmarshal(bytes, &config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}
