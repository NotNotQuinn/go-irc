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
