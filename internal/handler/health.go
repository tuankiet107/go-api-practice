package handler

import (
	"encoding/json"
	"net/http"
)

type healthResponse struct {
	Status string `json:"status"`
}

func Health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := healthResponse{Status: "ok"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// The response has already started, so there is no other HTTP response
		// that can safely be sent here.
		return
	}
}
