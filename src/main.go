package main

import (
	"fmt"
	"sync"

	"./bot"
	"./checker"
	"./config"
	"./jira"
	"./oauth"
	"./users"
)

func main() {
	fmt.Println("Main started!")

	config := config.New("config.yml")

	usersStorage := users.New()
	jira.Init(config)

	var wg sync.WaitGroup
	wg.Add(3)

	go oauth.Start(&wg, usersStorage, config)
	go bot.Start(&wg, usersStorage, config)
	go checker.Start(&wg, usersStorage, config)

	wg.Wait()
}
