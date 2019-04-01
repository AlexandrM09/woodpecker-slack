package slack

import (
	"fmt"

	"github.com/nlopes/slack"

	"../config"
)

type Client struct {
	api *slack.Client
}

type Message struct {
	User    string
	Channel ChannelID
	Text    string
}

type ChannelID string

func New(config *config.Config) *Client {
	fmt.Println("New config")
	fmt.Println(config.Slack.Token)
	a := slack.New(config.Slack.Token)
	return &Client{api: a}
}

func (c *Client) SendMessage(text string, Channel ChannelID) {
	c.api.PostMessage(Channel.GetRealID(), slack.MsgOptionText(text, false))
}

func (c *Client) GetMessages() chan Message {
	rtm := c.api.NewRTM()
	go rtm.ManageConnection()

	res := make(chan Message)
	go func() {
		for msg := range rtm.IncomingEvents {
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				if ev.SubType != "bot_message" {
					res <- Message{Text: ev.Text, User: ev.User, Channel: ChannelID(ev.Channel)}
				}
			}
		}
	}()

	return res
}

func (id ChannelID) GetRealID() string {
	return string(id)
}
