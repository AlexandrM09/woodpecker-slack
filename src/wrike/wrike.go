package wrike

import (
	"fmt"

	wrike "github.com/pierreboissinot/go-wrike"
)

type Client struct {
	api *wrike.Client
}

type users struct {
	Kind string
	Data []Data
}

type Data struct {
	ID        string
	FirstName string
	LastName  string
	Type      string
	Profiles  []profile
}

type profile struct {
	Email string
}

func filter(vs []Data, f func(Data) bool) []Data {
	vsf := make([]Data, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func New(token string) *Client {
	return &Client{api: wrike.NewClient(nil, token)}
}

func (c *Client) GetUsers() []Data {
	req, _ := c.api.NewRequest("GET", "contacts", nil)

	u := new(users)
	c.api.Do(req, u)
	d := filter(u.Data, func(d Data) bool {
		return d.Type == "Person" && d.LastName != "Bot"
	})

	return d
}

func GetUserIDByToken(token string) string {
	fmt.Println(token)
	api := wrike.NewClient(nil, token)
	req, _ := api.NewRequest("GET", "account", nil)

	u := new(struct {
		Data []struct {
			ID string
		}
	})
	api.Do(req, u)
	return u.Data[0].ID
}

func GetActiveTasksByUser() {

}

func GetPotentialTasksByUser() {

}
