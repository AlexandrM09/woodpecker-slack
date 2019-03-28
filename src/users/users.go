package users

import (
	"sync"
)

// MaxUsersCount means how many users can accommodate a users storage
const MaxUsersCount = 100

// SlackID is a representation of the type of user ID in slack
type SlackID string

// JiraID is a representation of the type of user ID in Jira
type JiraID string

// OauthToken is a representation of the type of user token in jira
type OauthToken string

// DuplicateError is an error thrown when you try to add an existing user
type DuplicateError struct {
	message string
}

func (e *DuplicateError) Error() string {
	return e.message
}

// Users is the storage of all known users.
type Users struct {
	mt    sync.RWMutex
	users []*User
}

// User is an abstraction over a user with accounts in jira and slack
type User struct {
	slackID    SlackID
	jiraID     JiraID
	oauthToken OauthToken
}

// New creates new users storage
func New() *Users {
	return &Users{users: make([]*User, 0, MaxUsersCount)}
}

// AddUser adds user to storage
func (users *Users) AddUser(user *User) error {
	users.mt.Lock()
	defer users.mt.Unlock()

	for _, u := range users.users {
		if u.slackID == user.slackID || u.jiraID == u.jiraID {
			return &DuplicateError{"User already exist"}
		}
	}
	users.users = append(users.users, user)
	return nil
}

// FindBySlackID finds user by slack id
func (users *Users) FindBySlackID(slackID SlackID) *User {
	users.mt.RLock()
	defer users.mt.RUnlock()

	for _, user := range users.users {
		if user.slackID == slackID {
			return user
		}
	}
	return nil
}

// FindByJiraID finds user by jira id
func (users *Users) FindByJiraID(jiraID JiraID) *User {
	users.mt.RLock()
	defer users.mt.RUnlock()

	for _, user := range users.users {
		if user.jiraID == jiraID {
			return user
		}
	}
	return nil
}
