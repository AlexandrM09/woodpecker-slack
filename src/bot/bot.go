package bot

import (
	"fmt"
	"sync"

	"../slack"
	musers "../users"
)

var api *slack.Client

// Start starts the bot
func Start(wg *sync.WaitGroup, users *musers.Users, api *slack.Client) {
	defer wg.Done()

	fmt.Println("Bot started!")

	for message := range api.GetMessages() {
		// Оставить комментарий к задаче
		// pattern := regexp.MustCompile("^Оставить комментарий \"([^\"]+)\" к задаче (\\S+)$")
		// Оставить затраченное время и комментарий
		// Взять задачу с комментарием
		// Перевести задачу в состояние готово
		// Перевести задачу на другого пользователя
		user := users.FindBySlackID(musers.SlackID(message.User))
		user.SlackChannal = string(message.Channel)
		if user.OauthToken == "" {
			api.SendMessage("Need oauth https://www.wrike.com/oauth2/authorize/v4?client_id=A61nSaq5&response_type=code", message.Channel)
		} else {
			api.SendMessage("Hello", message.Channel)
		}
	}
}
