package bot

import (
	"fmt"
	"sync"


	"../config"
	"../users"
	"../slack"
)

var api *slack.Client

// Start starts the bot
func Start(wg *sync.WaitGroup, users *users.Users, api *slack.Client, config *config.Config) {
	defer wg.Done()

	fmt.Println("Bot started!")

	for message := range api.GetMessages() {
		// Оставить комментарий к задаче
		// pattern := regexp.MustCompile("^Оставить комментарий \"([^\"]+)\" к задаче (\\S+)$")
		// Оставить затраченное время и комментарий
		// Взять задачу с комментарием
		// Перевести задачу в состояние готово
		// Перевести задачу на другого пользователя
		api.SendMessage("Hello", message.Channel)
	}
}
