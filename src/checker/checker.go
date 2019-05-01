package checker

import (
	"fmt"
	"sync"
	"time"

	"../slack"
	"../users"
	"../wrike"
)

// Start starts the checker
func Start(wg *sync.WaitGroup, users *users.Users, api *wrike.Client, apiM *slack.Client) {
	defer wg.Done()

	fmt.Println("Checker started!")

	date := time.Now().UTC()
	updateUsers(users, api, apiM)
	for _, user := range users.GetUsers() {
		fmt.Println(api.GetOutdatedTasksByUser(string(user.WrikeID), date))
	}
}

func updateUsers(us *users.Users, api *wrike.Client, apiM *slack.Client) {
	serverUsers := api.GetUsers()

	for _, user := range serverUsers {
		slackID, _ := apiM.GetIDByEmail(user.Profiles[0].Email)
		newUser := &users.User{WrikeID: users.WrikeID(user.ID), Email: user.Profiles[0].Email, SlackID: users.SlackID(slackID)}
		us.AddUser(newUser)
	}
}
