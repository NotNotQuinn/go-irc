package config

import (
	"encoding/json"
	"io/ioutil"
)

type PublicTwitchConfig struct {
	Channels []string `json:"channels"`
}

type PublicGlobalConfig struct {
	CommandPrefix string `json:"commandPrefix"`
}

type PublicConfig struct {
	Twitch PublicTwitchConfig `json:"twitch"`
	Global PublicGlobalConfig `json:"global"`
}

func GetPublic() (*PublicConfig, error) {
	bytes, err := ioutil.ReadFile("./config/public_conf.json")
	if err != nil {
		return nil, err
	}

	var config PublicConfig
	json.Unmarshal(bytes, &config)

	return &config, nil
}

// Handling private config

type PrivateConfig struct {
	// Even though the username is not 'private' its very much related so I will keep it here

	// Username of account
	Username string `json:"username"`
	// Oauth token of account
	Oauth string `json:"oauth"`
}

func GetPrivate() (conf *PrivateConfig, err error) {
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
