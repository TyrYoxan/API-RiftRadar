package http

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tyryoxan/API-RiftRadar/app/usecases"
	"net/http"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

type PlayerAction struct {
	Service *usecases.PlayerService
}

// respondWithJSON is a helper function to respond with JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		// If encoding fails, log the error and send a plain text response
		fmt.Printf("Error encoding response: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// respondWithError is a helper function to respond with an error
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, ErrorResponse{Error: message})
}

func (a *PlayerAction) ListPlayers(w http.ResponseWriter, req *http.Request) {
	// Simply list all players
	players := a.Service.GetPlayers()
	respondWithJSON(w, http.StatusOK, players)
}

func (a *PlayerAction) GetPlayerByNameAndRegion(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	name := vars["name"]
	region := vars["region"]

	if name == "" || region == "" {
		respondWithError(w, http.StatusBadRequest, "name and region parameters are required")
		return
	}

	player, err := a.Service.GetPlayerByNameAndRegion(name, region)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, player)
}

func (a *PlayerAction) GetPlayerByNameAndTag(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	name := vars["name"]
	tag := vars["tag"]

	if name == "" || tag == "" {
		respondWithError(w, http.StatusBadRequest, "name and tag parameters are required")
		return
	}

	player, err := a.Service.GetPlayerByNameAndTag(name, tag)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, player)
}
