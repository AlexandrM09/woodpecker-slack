package wrike

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	musers "../users"
	newWrike "github.com/DarkHole1/go-wrike"
	wrike "github.com/pierreboissinot/go-wrike"
)

type Client struct {
	api          *wrike.Client
	newAPI       *newWrike.API
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
	Admin bool
	Owner bool
}

var id, secret string

func filter(vs []Data, f func(Data) bool) []Data {
	vsf := make([]Data, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func New(token, _id, _secret string) *Client {
	id = _id
	secret = _secret
	client := &Client{
		api:          wrike.NewClient(nil, token),
		newAPI:       &newWrike.API{Token: token, ID: _id, Secret: _secret},
		nameToStatus: make(map[string]string),
		statusToName: make(map[string]string),
	}

	workflows, _ := client.newAPI.GetWorkflows()

	for _, workflow := range workflows {
		if workflow.Standard {
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
	fmt.Println("Token: " + token)
	api := wrike.NewClient(nil, token)
	var data struct {
		Me bool `url:"me"`
	}
	data.Me = true
	req, _ := api.NewRequest("GET", "contacts", data)

	u := new(struct {
		Data []struct {
			ID string
		}
	})
	api.Do(req, u)
	fmt.Println(u)
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
	var params struct {
		Responsibles   string `url:"responsibles"`
		CustomStatuses string `url:"customStatuses"`
	}
	params.Responsibles = "[" + id + "]"
	params.CustomStatuses = "[" + c.nameToStatus["New"] + "]"

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

func (c *Client) CommentTask(id, comment string) (bool, error) {
	req, err := c.api.NewRequest("POST", "tasks/"+id+"/comments", nil)
	if err != nil {
		return false, err
	}
	form := url.Values{"text": {comment}, "plainText": {"true"}}.Encode()
	req.Body = ioutil.NopCloser(strings.NewReader(form))
	req.ContentLength = int64(len(form))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp := new(struct {
		Error            string
		ErrorDescription string
	})

	_, err = c.api.Do(req, resp)
	if err != nil {
		return false, err
	}

	if resp.Error != "" {
		fmt.Println(resp)
		return false, errors.New(resp.ErrorDescription)
	}

	return true, nil
}

func sfilter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func (c *Client) TakeTask(userid, taskid string) (bool, error) {
	// query task
	req1, err := c.api.NewRequest("GET", "tasks/"+taskid, nil)
	if err != nil {
		return false, err
	}
	resp1 := new(struct {
		Error            string
		ErrorDescription string
		Data             []struct {
			ResponsibleIDs []string
		}
	})
	_, err = c.api.Do(req1, resp1)
	if err != nil {
		return false, err
	}

	if resp1.Error != "" {
		return false, errors.New(resp1.ErrorDescription)
	}

	needToDelete := resp1.Data[0].ResponsibleIDs
	needToDelete = sfilter(needToDelete, func(s string) bool {
		return s != userid
	})
	// remove responsibles and change status
	req2, err := c.api.NewRequest("PUT", "tasks/"+taskid, nil)
	if err != nil {
		return false, err
	}

	form := url.Values{"removeResponsibles": {"[" + strings.Join(needToDelete, ", ") + "]"}, "customStatus": {c.nameToStatus["In Progress"]}}.Encode()
	req2.Body = ioutil.NopCloser(strings.NewReader(form))
	req2.ContentLength = int64(len(form))
	req2.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	fmt.Println(form)

	resp2 := new(struct {
		Error            string
		ErrorDescription string
	})

	_, err = c.api.Do(req2, resp2)
	if err != nil {
		return false, err
	}

	if resp2.Error != "" {
		return false, errors.New(resp2.ErrorDescription)
	}

	return true, nil
}

func (c *Client) FinishTask(taskid string) (bool, error) {
	req, err := c.api.NewRequest("PUT", "tasks/"+taskid, nil)
	if err != nil {
		return false, err
	}

	form := url.Values{"customStatus": {c.nameToStatus["Completed"]}}.Encode()
	req.Body = ioutil.NopCloser(strings.NewReader(form))
	req.ContentLength = int64(len(form))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	fmt.Println(form)

	resp := new(struct {
		Error            string
		ErrorDescription string
	})

	_, err = c.api.Do(req, resp)
	if err != nil {
		return false, err
	}

	if resp.Error != "" {
		return false, errors.New(resp.ErrorDescription)
	}

	return true, nil
}

func (c *Client) MoveTask(taskid, userid string) (bool, error) {
	// query task
	req1, err := c.api.NewRequest("GET", "tasks/"+taskid, nil)
	if err != nil {
		return false, err
	}
	resp1 := new(struct {
		Error            string
		ErrorDescription string
		Data             []struct {
			ResponsibleIDs []string
		}
	})
	_, err = c.api.Do(req1, resp1)
	if err != nil {
		return false, err
	}

	if resp1.Error != "" {
		return false, errors.New(resp1.ErrorDescription)
	}

	needToDelete := resp1.Data[0].ResponsibleIDs

	// remove responsibles and change status
	req2, err := c.api.NewRequest("PUT", "tasks/"+taskid, nil)
	if err != nil {
		return false, err
	}

	form := url.Values{"removeResponsibles": {"[" + strings.Join(needToDelete, ", ") + "]"}, "addResponsibles": {"[" + userid + "]"}}.Encode()
	req2.Body = ioutil.NopCloser(strings.NewReader(form))
	req2.ContentLength = int64(len(form))
	req2.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	fmt.Println(form)

	resp2 := new(struct {
		Error            string
		ErrorDescription string
	})

	_, err = c.api.Do(req2, resp2)
	if err != nil {
		return false, err
	}

	if resp2.Error != "" {
		return false, errors.New(resp2.ErrorDescription)
	}

	return true, nil
}

func (c *Client) GetOutlastedTasksWithoutUser(date time.Time) []Task {
	var params taskParams
	params.Responsibles = "[]"
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

func (c *Client) GetVeryOutdatedTasks(date time.Time) []Task {
	var params struct {
		CustomStatuses string `url:"customStatuses"`
		UpdatedDate    string `url:"updatedDate"`
	}
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

type Project struct {
	ID, Title string
}

func (c *Client) GetProjects() []Project {
	var params struct {
		Project bool `url:"project"`
		Deleted bool `url:"deleted"`
	}

	params.Project = true
	params.Deleted = false

	req, _ := c.api.NewRequest("GET", "/folders", params)
	resp := new(struct {
		Data []Project
	})

	_, err := c.api.Do(req, resp)
	if err != nil {
		panic(err)
	}

	return resp.Data
}

func (c *Client) FromOAuth(user *musers.User) *Client {
	access := string(user.OauthToken)
	refresh := user.RefreshToken
	api := &Client{nameToStatus: c.nameToStatus, statusToName: c.statusToName}
	api.api = wrike.NewClient(nil, access)

	if !api.Check() {
		access, refresh = Refresh(refresh)
		user.OauthToken = musers.OauthToken(access)
		user.RefreshToken = refresh
		api.api = wrike.NewClient(nil, access)
		if !api.Check() {
			return nil
		}
	}

	return api
}

func (c *Client) Check() bool {
	req, _ := c.api.NewRequest("GET", "version", nil)

	resp := new(struct {
		Error string
	})
	c.api.Do(req, resp)

	return len(resp.Error) == 0
}

func Refresh(refresh string) (string, string) {
	resp, _ := http.PostForm("https://www.wrike.com/oauth2/token", url.Values{
		"client_id":     {id},
		"client_secret": {secret},
		"grant_type":    {"refresh_token"},
		"refresh_token": {refresh},
	})

	byteBody, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var body map[string]interface{}
	json.Unmarshal(byteBody, &body)

	if val, ok := body["error"].(string); ok {
		fmt.Println(val)
		fmt.Println("Error: " + body["error_description"].(string))
		return "", ""
	} else {
		return body["access_token"].(string), body["refresh_token"].(string)
	}
}
