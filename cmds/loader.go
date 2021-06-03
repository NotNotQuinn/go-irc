package cmd

// seperate file for

// Loads a all loads all commands
func LoadAll() {
	pingCommand.load()
	commandCommand.load()
	aboutCommand.load()
	githubCommand.load()
}
