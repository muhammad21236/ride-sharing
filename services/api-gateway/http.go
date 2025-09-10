package main

import (
	"encoding/json"
	"net/http"
)

func HandleTripPreview(w http.ResponseWriter, r *http.Request) {
	var reqBody previewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if reqBody.UserID == "" {
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"status": "preview received"})
}
