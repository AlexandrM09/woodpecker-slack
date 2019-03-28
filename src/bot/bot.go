package bot

import (
	"fmt"
	"sync"

	"github.com/nlopes/slack"

	"../config"
	"../users"
)

var api *slack.Client

// Start starts the bot
func Start(wg *sync.WaitGroup, users *users.Users, config *config.Config) {
	defer wg.Done()

	fmt.Println("Bot started!")

	api = slack.New(config.Slack.Token)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			if ev.SubType != "bot_message" {
				messageRouter(ev.Text, ev.User, ev.Channel)
			}
		}
	}
}

func messageRouter(text, user, channel string) {
	// Оставить комментарий к задаче
	// pattern := regexp.MustCompile("^Оставить комментарий \"([^\"]+)\" к задаче (\\S+)$")
	// Оставить затраченное время и комментарий
	// Взять задачу с комментарием
	// Перевести задачу в состояние готово
	// Перевести задачу на другого пользователя
	api.PostMessage(channel, slack.MsgOptionText("Hello", false))
}
