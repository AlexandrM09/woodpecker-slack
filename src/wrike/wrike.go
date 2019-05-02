package wrike

import (
	"fmt"
	"time"

	wrike "github.com/pierreboissinot/go-wrike"
)

type Client struct {
	api          *wrike.Client
	statusToName map[string]string
	nameToStatus map[string]string
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
	client := &Client{
		api:          wrike.NewClient(nil, token),
		nameToStatus: make(map[string]string),
		statusToName: make(map[string]string),
	}

	data := new(struct {
		Data []struct {
			Name           string
			CustomStatuses []struct {
				ID   string
				Name string
			}
		}
	})

	req, _ := client.api.NewRequest("GET", "workflows", nil)
	client.api.Do(req, data)

	for _, workflow := range data.Data {
		if workflow.Name == "Default Workflow" {
			for _, status := range workflow.CustomStatuses {
				client.nameToStatus[status.Name] = status.ID
				client.statusToName[status.ID] = status.Name
			}
		}
	}

	return client
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

type Task struct {
	ID             string
	Title          string
	CustomStatus   string
	CustomStatusID string
	UpdatedDate    string
}

type taskParams struct {
	Responsibles   string `url:"responsibles"`
	CustomStatuses string `url:"customStatuses"`
	UpdatedDate    string `url:"updatedDate"`
}

type tasksResponse struct {
	Data []Task
}

func (c *Client) GetOutdatedTasksByUser(id string, date time.Time) []Task {
	var params taskParams
	params.Responsibles = "[" + id + "]"
	params.CustomStatuses = "[" + c.nameToStatus["In Progress"] + "]"
	params.UpdatedDate = "{\"end\":\"" + date.UTC().Format("2006-01-02T15:04:05Z") + "\"}"

	req, _ := c.api.NewRequest("GET", "tasks", params)
	resp := new(tasksResponse)
	_, err := c.api.Do(req, resp)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(resp.Data); i++ {
		resp.Data[i].CustomStatus = c.statusToName[resp.Data[i].CustomStatusID]
	}

	return resp.Data
}

func (c *Client) GetTasksInProgressByUser(id string) []Task {
	return c.GetOutdatedTasksByUser(id, time.Now())
}

func (c *Client) GetPotentialTasksByUser(id string) []Task {
	var params taskParams
	params.Responsibles = "[" + id + "]"
	params.CustomStatuses = "[" + c.nameToStatus["New"] + "]"
	params.UpdatedDate = "{}"

	req, _ := c.api.NewRequest("GET", "tasks", params)
	resp := new(tasksResponse)
	_, err := c.api.Do(req, resp)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(resp.Data); i++ {
		resp.Data[i].CustomStatus = c.statusToName[resp.Data[i].CustomStatusID]
	}

	return resp.Data
}
