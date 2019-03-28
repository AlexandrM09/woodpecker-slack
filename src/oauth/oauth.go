package oauth

import (
	"fmt"
	"sync"

	"../config"
	"../users"
)

// Start starts the oauth server
func Start(wg *sync.WaitGroup, users *users.Users, config *config.Config) {
	defer wg.Done()

	fmt.Println("Oauth server started!")
}
