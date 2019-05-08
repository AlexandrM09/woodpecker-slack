package users

import (
	"sync"

	bolt "github.com/boltdb/bolt"
)

// MaxUsersCount means how many users can accommodate a users storage
const MaxUsersCount = 100

// SlackID is a representation of the type of user ID in slack
type SlackID string

// WrikeID is a representation of the type of user ID in Wrike
type WrikeID string

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
	mt     sync.RWMutex
	dbFile string
	db     *bolt.DB
	users  []*User
}

// User is an abstraction over a user with accounts in jira and slack
type User struct {
	SlackID      SlackID
	WrikeID      WrikeID
	OauthToken   OauthToken
	RefreshToken string
	SlackChannal string
	Email        string
}

// New creates new users storage
func New(dbFile string) *Users {
	users := &Users{users: make([]*User, 0, MaxUsersCount), dbFile: dbFile}

	if dbFile != "" {
		db, err := bolt.Open(dbFile, 0600, nil)
		if err != nil {
			return nil
		}
		users.db = db
	}

	return users
}

// AddUser adds user to storage
func (users *Users) AddUser(user *User) error {
	users.mt.Lock()
	defer users.mt.Unlock()

	for _, u := range users.users {
		if u.SlackID == user.SlackID || u.WrikeID == u.WrikeID {
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
		if user.SlackID == slackID {
			return user
		}
	}
	return nil
}

// FindByWrikeID finds user by wrike id
func (users *Users) FindByWrikeID(wrikeID WrikeID) *User {
	users.mt.RLock()
	defer users.mt.RUnlock()

	for _, user := range users.users {
		if user.WrikeID == wrikeID {
			return user
		}
	}
	return nil
}

func (users *Users) GetUsers() []*User {
	users.mt.RLock()
	defer users.mt.RUnlock()

	tmp := make([]*User, len(users.users))
	copy(tmp, users.users)

	return tmp
}

func (users *Users) Sync() error {

}
