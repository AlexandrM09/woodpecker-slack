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

	updateUsers(users, api, apiM)

	for {
		// time.Sleep(15 * time.Minute)
		time.Sleep(5 * time.Second)

		date := time.Now()
		if !checkWeekends(date) {
			for {
				date = date.AddDate(0, 0, -1)
				if !checkWeekends(date) {
					break
				}
			}

			for _, user := range users.GetUsers() {
				if user.SlackChannal == "" {
					continue
				}
				go processUser(user, date, api, apiM)
			}
		}
	}
}

func updateUsers(us *users.Users, api *wrike.Client, apiM *slack.Client) {
	serverUsers := api.GetUsers()

	for _, user := range serverUsers {
		slackID, _ := apiM.GetIDByEmail(user.Profiles[0].Email)
		newUser := &users.User{WrikeID: users.WrikeID(user.ID), Email: user.Profiles[0].Email, SlackID: users.SlackID(slackID)}
		fmt.Println(newUser)
		us.AddUserIfNotExist(newUser)
	}
}

func checkWeekends(date time.Time) bool {
	return false
	// weekday := date.Weekday()
	// if weekday == 0 || weekday == 6 {
	// 	return true
	// }
	// return false
}

func processUser(user *users.User, date time.Time, api *wrike.Client, apiM *slack.Client) {
	tasks := api.GetTasksInProgressByUser(string(user.WrikeID))
	if len(tasks) != 0 {
		tasks = api.GetOutdatedTasksByUser(string(user.WrikeID), date)
		if len(tasks) != 0 {
			s := "You have some outdated tasks: \n"
			for _, task := range tasks {
				s += "- " + task.ID + ": " + task.Title + "\n"
			}
			apiM.SendMessage(s, slack.ChannelID(user.SlackChannal))
		} else {
			apiM.SendMessage("Everything is ok", slack.ChannelID(user.SlackChannal))
		}
	} else {
		tasks = api.GetPotentialTasksByUser(string(user.WrikeID))
		if len(tasks) != 0 {
			s := "You don't have any tasks in progress. Choose one:\n"
			for _, task := range tasks {
				s += "- " + task.ID + ": " + task.Title + "\n"
			}
			apiM.SendMessage(s, slack.ChannelID(user.SlackChannal))
		} else {
			apiM.SendMessage("You don't have any tasks ;(", slack.ChannelID(user.SlackChannal))
		}
	}
}
