package users

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
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
	gob.Register(User{})
	users := &Users{users: make([]*User, 0, MaxUsersCount), dbFile: dbFile}

	if dbFile != "" {
		db, err := bolt.Open(dbFile, 0600, nil)
		if err != nil {
			return nil
		}
		users.db = db
		err = users.Load()
		if err != nil {
			return nil
		}
	}

	return users
}

// AddUser adds user to storage
func (users *Users) AddUser(user *User, sync bool) error {
	users.mt.Lock()

	for _, u := range users.users {
		if u.SlackID == user.SlackID || u.WrikeID == user.WrikeID {
			users.mt.Unlock()
			return &DuplicateError{"User already exist"}
		}
	}
	users.users = append(users.users, user)

	users.mt.Unlock()
	if sync {
		return users.Sync()
	}
	return nil
}

func (users *Users) AddUserIfNotExist(user *User) {
	users.AddUser(user, true)
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
	if users.db == nil {
		return errors.New("No database")
	}
	users.mt.RLock()
	defer users.mt.RUnlock()

	err := users.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return err
		}

		for _, user := range users.users {
			b := bytes.Buffer{}
			e := gob.NewEncoder(&b)
			err := e.Encode(user)
			if err != nil {
				return err
			}

			// fmt.Println(b.Bytes())
			bucket.Put([]byte(user.WrikeID), b.Bytes())
		}
		return nil
	})

	return err
}

func (users *Users) Close() {
	users.db.Close()
}

func (users *Users) Load() error {
	if users.db == nil {
		return errors.New("No database")
	}

	users.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		if bucket == nil {
			return errors.New("No bucket")
		}

		err := bucket.ForEach(func(k []byte, v []byte) error {
			var user User
			b := bytes.Buffer{}
			b.Write(v)
			d := gob.NewDecoder(&b)
			err := d.Decode(&user)
			if err != nil {
				return err
			}
			users.AddUser(&user, false)
			return nil
		})

		return err
	})

	fmt.Println("Here")
	fmt.Println(users.users)
	return nil
}
