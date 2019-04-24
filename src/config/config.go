package config

import (
	"fmt"
	"os"

	cfg "github.com/olebedev/config"
)

// Config is config
type Config struct {
	Slack slack
	Wrike wrike
}

type slack struct {
	Token string
}

type jira struct {
	Username string
	Password string
}

type wrike struct {
	Token  string
	ID     string
	Secret string
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

	config.Wrike.Token, err = configReader.String("wrike.token")
	if err != nil {
		return nil
	}

	config.Wrike.ID, err = configReader.String("wrike.id")
	if err != nil {
		return nil
	}

	config.Wrike.Secret, err = configReader.String("wrike.secret")
	if err != nil {
		return nil
	}

	return &config
}
