package http

import (
	"encoding/json"
	"github.com/tyryoxan/API-RiftRadar/app/usecases"
	"io"
	"net/http"
)

type LeaderboardAction struct {
	Service *usecases.LeaderboardService
}

type payload struct {
	Name string `json:"name"`
	Tag  string `json:"tag"`
}

func (a *LeaderboardAction) GetLeaderboard(w http.ResponseWriter, req *http.Request) {
	leaderboard := a.Service.GetLeaderboard()
	respondWithJSON(w, http.StatusOK, leaderboard)
}

func (a *LeaderboardAction) GetLeaderboardByRegion(w http.ResponseWriter, req *http.Request) {
	region := req.URL.Query().Get("region")
	if region == "" {
		respondWithError(w, http.StatusBadRequest, "Region parameter is required")
		return
	}

	leaderboard, err := a.Service.GetLeaderboardByRegion(region)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, leaderboard)
}

func (a *LeaderboardAction) CreateLeaderboard(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer req.Body.Close()

	var name payload
	json.Unmarshal(body, &name)

	learderboard := a.Service.CreateLeaderboard(name.Name)
	respondWithJSON(w, http.StatusOK, learderboard)
}

func (a *LeaderboardAction) AddPlayer(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer req.Body.Close()
	var player payload
	json.Unmarshal(body, &player)

	err = a.Service.AddPlayer(player.Name, player.Tag)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Player added successfully"})
}
