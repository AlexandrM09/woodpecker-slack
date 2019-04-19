package oauth

import (
	"fmt"
	"net/http"
	"sync"

	"../users"
)

// Start starts the oauth server
func Start(wg *sync.WaitGroup, users *users.Users) {
	defer wg.Done()

	fmt.Println("Oauth server started!")

	http.HandleFunc("/", oauthHandler)
	http.ListenAndServe(":9000", nil)
}

func oauthHandler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["code"]
	if !ok || len(keys[0]) < 1 {
		fmt.Fprintln(w, "Not ok")
	} else {
		fmt.Println(keys[0])
		fmt.Fprintln(w, "Hello")
	}
}
