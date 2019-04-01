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
)

func main() {
	fmt.Println("Main started!")

	config := config.New("config.yml")
	if config == nil {
		panic("Config didn't load")
	}
	api := slack.New(config)

	usersStorage := users.New()

	var wg sync.WaitGroup
	wg.Add(3)

	go oauth.Start(&wg, usersStorage, config)
	go bot.Start(&wg, usersStorage, api, config)
	go checker.Start(&wg, usersStorage, config)

	wg.Wait()
}
