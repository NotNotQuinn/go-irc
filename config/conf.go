package config

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
)

var confFile = "./config/public_conf.json"
var Public *PublicConfig
var Private *PrivateConfig

type PublicTwitchConfig struct {
	Channels []string `json:"channels"`
}

type PublicGlobalConfig struct {
	CommandPrefix  string `json:"commandPrefix"`
	Admin_Username string `json:"admin_username"`
	UserAgent      string `json:"user_agent"`
}

type PublicConfig struct {
	Twitch PublicTwitchConfig `json:"twitch"`
	Global PublicGlobalConfig `json:"global"`
}

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

func getPublic() (*PublicConfig, error) {
	bytes, err := ioutil.ReadFile(confFile)
	if err != nil {
		return nil, err
	}
	var config PublicConfig
	json.Unmarshal(bytes, &config)
	return &config, nil
}

func (conf *PublicConfig) Reload() error {
	bytes, err := ioutil.ReadFile(confFile)
	if err != nil {
		return err
	}
	var config PublicConfig
	json.Unmarshal(bytes, &config)
	return nil
}

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

// Handling private config

type PrivateConfig struct {
	// Even though the username is not 'private' its very much related so I will keep it here

	// Username of account
	Username string `json:"username"`
	// Oauth token of account
	Oauth string `json:"oauth"`
}

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
