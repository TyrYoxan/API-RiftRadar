package routes

import (
	"github.com/gorilla/mux"
	httpAdapter "github.com/tyryoxan/API-RiftRadar/adapters/in/http"
)

// SetupRouter configures all the routes for the application and returns a router
func SetupRouter(playerAction *httpAdapter.PlayerAction, leaderboardAction *httpAdapter.LeaderboardAction) *mux.Router {
	r := mux.NewRouter()

	// Players routes
	r.HandleFunc("/players", playerAction.ListPlayers).Methods("GET")
	r.HandleFunc("/players/{name}/{tag}", playerAction.GetPlayerByNameAndTag).Methods("GET")

	// Leaderboard routes
	r.HandleFunc("/leaderboard", leaderboardAction.GetLeaderboard).Methods("GET")
	r.HandleFunc("/leaderboard", leaderboardAction.CreateLeaderboard).Methods("POST")
	r.HandleFunc("/leaderboard/player", leaderboardAction.AddPlayer).Methods("POST")

	return r
}
