package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Username string `json:"username"`
	Oauth    string `json:"oauth"`
}

func Get() (conf *Config, err error) {
	bytes, err := ioutil.ReadFile("./config/conf.json")
	if err != nil {
		return nil, err
	}

	var config Config
	json.Unmarshal(bytes, &config)

	return &config, nil
}
