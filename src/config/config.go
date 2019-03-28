package config

import (
	"fmt"
	"os"

	cfg "github.com/olebedev/config"
)

// Config is config
type Config struct {
	Slack slack
	Jira  jira
}

type slack struct {
	Token string
}

type jira struct {
	Username string
	Password string
}

// New loads config from filename
func New(filename string) *Config {
	fmt.Println("Loading config")

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil
	}

	var err error
	var configReader *cfg.Config
	configReader, err = cfg.ParseYamlFile(filename)

	if err != nil {
		return nil
	}

	config := Config{}

	config.Slack.Token, err = configReader.String("slack.token")
	if err != nil {
		return nil
	}

	config.Jira.Username, err = configReader.String("jira.username")
	if err != nil {
		return nil
	}

	config.Jira.Password, err = configReader.String("jira.password")
	if err != nil {
		return nil
	}

	return &config
}
