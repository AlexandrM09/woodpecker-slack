package main

import (
	"fmt"
	"sync"

	"./bot"
	"./checker"
	"./config"
	"./oauth"
	"./slack"
	"./users"
	"./wrike"
)

func main() {
	fmt.Println("Main started!")

	config := config.New("config.yml")
	if config == nil {
		panic("Config didn't load")
	}

	apiMessenger := slack.New(config.Slack.Token)
	apiTaskmanager := wrike.New(config.Wrike.Token)
	usersStorage := users.New("woodpecker.db")
	defer usersStorage.Close()

	var wg sync.WaitGroup
	wg.Add(3)

	go oauth.Start(&wg, usersStorage, config)
	go bot.Start(&wg, usersStorage, apiTaskmanager, apiMessenger, config)
	go checker.Start(&wg, usersStorage, apiTaskmanager, apiMessenger)

	wg.Wait()
}
