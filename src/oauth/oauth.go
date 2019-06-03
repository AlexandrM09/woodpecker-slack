package oauth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"

	"../config"
	"../users"
	"../wrike"
)

var id string
var secret string
var localUsers *users.Users

// Start starts the oauth server
func Start(wg *sync.WaitGroup, users *users.Users, config *config.Config) {
	defer wg.Done()

	fmt.Println("Oauth server started!")

	id = config.Wrike.ID
	secret = config.Wrike.Secret
	localUsers = users

	http.HandleFunc("/", oauthHandler)
	http.ListenAndServe(":9000", nil)
}

func oauthHandler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["code"]
	if !ok || len(keys[0]) < 1 {
		fmt.Fprintln(w, "Not ok")
	} else {
		fmt.Println(keys[0])
		resp, err := http.PostForm("https://www.wrike.com/oauth2/token", url.Values{
			"client_id":     {id},
			"client_secret": {secret},
			"grant_type":    {"authorization_code"},
			"code":          {keys[0]},
		})

		if err != nil {
			fmt.Println(err)
		} else {
			byteBody, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()

			var body map[string]interface{}
			err := json.Unmarshal(byteBody, &body)

			if err != nil {
				fmt.Println(err)
			}

			if val, ok := body["error"].(string); ok {
				fmt.Println(val)
				fmt.Println("Error: " + body["error_description"].(string))
			} else {
				// fmt.Println(body)
				accessToken := body["access_token"].(string)
				refreshToken := body["refresh_token"].(string)
				user := localUsers.FindByWrikeID(users.WrikeID(wrike.GetUserIDByToken(accessToken)))
				fmt.Println(user)
				if user == nil {
					fmt.Fprintln(w, "I don't know you")
				} else {
					user.OauthToken = users.OauthToken(accessToken)
					user.RefreshToken = refreshToken
					localUsers.Sync()
					fmt.Fprintln(w, "Success")
				}
			}

			// fmt.Println(body)
		}
	}
}
