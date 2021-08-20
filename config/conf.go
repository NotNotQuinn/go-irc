package config

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

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
	Twitch      PublicTwitchConfig `json:"twitch"`
	Global      PublicGlobalConfig `json:"global"`
	Users       PublicUsersConfig  `json:"users"`
	Development struct {
		Channels []string `json:"channels"`
		Prefix   string   `json:"prefix"`
	} `json:"development"`
	Production bool
	// The file the config was loaded from
	originFile string
}

// Load the configs and assign them to the variables
func init() {
	// Directory to look for config files
	var confDir = "./config"
	var privConfPath = filepath.Join(confDir, "private_conf.json")

	fmt.Println("WB_TEST = " + os.Getenv("WB_TEST"))
	if os.Getenv("WB_TEST") == "true" {
		// When inside docker container
		confDir = "/bot/config/"
		privConfPath = filepath.Join(confDir, "tests_private_conf.json")
	}

	var confFile = filepath.Join(confDir, "public_conf.json")
	var err error
	Private, err = getPrivate(privConfPath)
	if err != nil {
		panic(err)
	}
	Public, err = getPublic(confFile)
	if err != nil {
		panic(err)
	}
}

// Load and return the public config
func getPublic(confFile string) (*PublicConfig, error) {
	path, err := filepath.Abs(confFile)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config PublicConfig
	err = json.Unmarshal(bytes, &config)
	if _, err := os.Stat(filepath.Join(filepath.Join(confFile, ".."), "PRODUCTION")); err == nil {
		// Only set on load, not reload.
		config.Production = true
	} else if os.IsNotExist(err) {
		// Development
		config.Twitch.Channels = config.Development.Channels
		config.Global.CommandPrefix = config.Development.Prefix
	}
	return &config, err
}

// Reload the config from file
func (conf *PublicConfig) Reload() error {
	path, err := filepath.Abs(conf.originFile)
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	var config *PublicConfig
	err = json.Unmarshal(bytes, config)
	if err != nil {
		return err
	}
	if !conf.Production {
		// Development
		config.Twitch.Channels = config.Development.Channels
		config.Global.CommandPrefix = config.Development.Prefix
	}
	*conf = *config
	return nil
}

// Save the config to the file
func (conf *PublicConfig) Save() (success bool, err error) {
	bytes, err := json.MarshalIndent(conf, "", "\t")
	if err != nil {
		return false, err
	}
	path, err := filepath.Abs(conf.originFile)
	if err != nil {
		return false, err
	}
	file, err := os.Create(path)
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
	err = ioutil.WriteFile(path, bytes, fs.ModeType)
	return true, err
}

// All private config data
type PrivateConfig struct {
	// Even though the username is not 'private' its very much related so I will keep it here

	// Username of account
	Username string `json:"username"`
	// Oauth token of account
	Oauth string `json:"oauth"`
	// Database config
	Database PrivateDatabaseConfig `json:"database"`
}

// All private database config data
type PrivateDatabaseConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

// Creates a mariadb driver specific string to connect to the database on a specific database.
func (D *PrivateDatabaseConfig) ConnecterString(database string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", D.Username, D.Password, D.Host, D.Port, database)
}

// Load and return private config
func getPrivate(privConfPath string) (conf *PrivateConfig, err error) {
	path, err := filepath.Abs(privConfPath)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config PrivateConfig
	err = json.Unmarshal(bytes, &config)

	if err != nil {
		return nil, err
	}

	if os.Getenv("WB_TEST") == "true" {
		// Load test config and overwite database connection information.
		testsConfFile, err := filepath.Abs(filepath.Join(filepath.Join(privConfPath, ".."), "tests_private_conf.json"))
		if err != nil {
			return nil, err
		}
		bytes, err := ioutil.ReadFile(testsConfFile)
		if err != nil {
			return nil, err
		}

		var testsConfig PrivateConfig
		err = json.Unmarshal(bytes, &testsConfig)
		if err != nil {
			return nil, err
		}
		config.Database = testsConfig.Database
	}

	return &config, nil
}
