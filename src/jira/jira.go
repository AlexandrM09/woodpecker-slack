package jira

import (
	"fmt"

	jira "github.com/andygrunwald/go-jira"

	"../config"
)

// Client for jira's api
type Client struct {
	originalClient *jira.Client
}

// JiraClient is global client instance
var JiraClient *Client

// Init the api client
func Init(config *config.Config) error {
	fmt.Println("Hello")

	tp := jira.BasicAuthTransport{
		Username: config.Jira.Username,
		Password: config.Jira.Password,
	}

	jiraClient, err := jira.NewClient(tp.Client(), "http://woodpecker-test.atlassian.net")

	if err != nil {
		return err
	}

	JiraClient = &Client{jiraClient}

	return nil
}

// GetAllUsers gets all jira's users
func (client *Client) GetAllUsers() {

}

func GetAllActiveIssuesByUser() {

}
