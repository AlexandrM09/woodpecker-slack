package bot

import (
	"fmt"
	"regexp"
	"sync"

	"../config"
	"../slack"
	musers "../users"
)

var api *slack.Client

// Start starts the bot
func Start(wg *sync.WaitGroup, users *musers.Users, api *slack.Client, config *config.Config) {
	defer wg.Done()

	fmt.Println("Bot started!")

	// Оставить комментарий к задаче
	pattern := regexp.MustCompile(`^Оставить комментарий "([^"]+)" к задаче (\S+)$`)
	// Оставить затраченное время и комментарий
	pattern2 := regexp.MustCompile(`^Оставить затраченное время (\d+) часов и комментарий "([^"]+)" к задаче (\S+)$`)
	// Взять задачу с комментарием
	pattern3 := regexp.MustCompile(`^Взять задачу (\S+) с комментарием "([^"]+)"$`)
	// Перевести задачу в состояние готово
	pattern4 := regexp.MustCompile(`^Перевести задачу (\S+) в состояние готово$`)
	// Перевести задачу на другого пользователя
	pattern5 := regexp.MustCompile(`^Перевести задачу (\S+) на пользователя (\S+)$`)

	for message := range api.GetMessages() {
		user := users.FindBySlackID(musers.SlackID(message.User))
		user.SlackChannal = string(message.Channel)
		// if user.OauthToken == "" {
		if false {
			api.SendMessage("Need oauth https://www.wrike.com/oauth2/authorize/v4?client_id="+config.Wrike.ID+"&response_type=code", message.Channel)
		} else {
			if match := pattern.FindStringSubmatch(message.Text); match != nil {
				api.SendMessage("Comment on task "+match[2]+" left", message.Channel)

			} else if match := pattern2.FindStringSubmatch(message.Text); match != nil {
				api.SendMessage("Comment on task "+match[3]+" left", message.Channel)

			} else if match := pattern3.FindStringSubmatch(message.Text); match != nil {
				api.SendMessage("Take task "+match[1], message.Channel)

			} else if match := pattern4.FindStringSubmatch(message.Text); match != nil {
				api.SendMessage("Task "+match[1]+" finished", message.Channel)

			} else if match := pattern5.FindStringSubmatch(message.Text); match != nil {
				api.SendMessage("Moved task "+match[1]+" on user "+match[2], message.Channel)

			} else {
				api.SendMessage("Unrecongonized command", message.Channel)
			}
		}
	}
}
