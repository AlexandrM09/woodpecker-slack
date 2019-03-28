package checker

import (
	"fmt"
	"sync"

	"../config"
	"../users"
)

// Start starts the checker
func Start(wg *sync.WaitGroup, users *users.Users, config *config.Config) {
	defer wg.Done()

	fmt.Println("Checker started!")
}
