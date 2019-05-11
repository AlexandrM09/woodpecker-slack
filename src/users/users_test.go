package users

import "testing"

func TestNew(t *testing.T) {
	users := New("")

	if users == nil {
		t.Error("New created nil object")
		return
	}

	if cap(users.users) != MaxUsersCount {
		t.Error("Users capacity does not match MaxUsersCount")
		return
	}
}

func TestAddUser(t *testing.T) {
	users := New("")

	if len(users.users) != 0 {
		t.Error("Was created not clear storage")
		return
	}

	user := User{}
	err := users.AddUser(&user)

	if err != nil {
		t.Error("An error occured while normal adding")
		return
	}

	if len(users.users) == 0 {
		t.Error("User was not added")
		return
	}

	if users.users[0] != &user {
		t.Error("Corrupted user")
		return
	}

	err = users.AddUser(&user)
	if err == nil {
		t.Error("Successfully added duplicate")
	}
}

func TestFindBySlackID(t *testing.T) {
	const slackID = "U0G9QF9C6"
	users := New("")

	if users.FindBySlackID(slackID) != nil {
		t.Error("Found non-existent user")
		return
	}

	user := User{SlackID: slackID}
	users.AddUser(&user)

	if users.FindBySlackID(slackID) == nil {
		t.Error("Not found existent user")
		return
	}
}

func TestFindByJiraID(t *testing.T) {
	const wrikeID = "5ac5326a95d30150501e5ff4"
	users := New("")

	if users.FindByWrikeID(wrikeID) != nil {
		t.Error("Found non-existent user")
		return
	}

	user := User{WrikeID: wrikeID}
	users.AddUser(&user)

	if users.FindByWrikeID(wrikeID) == nil {
		t.Error("Not found existent user")
		return
	}
}

func TestError(t *testing.T) {
	const str = "randomstr"
	e := &DuplicateError{str}
	if e.Error() != str {
		t.Error("Duplicate error returns incorrect message")
	}
}
