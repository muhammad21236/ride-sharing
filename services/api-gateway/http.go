package main

import (
	"encoding/json"
	"net/http"
	"ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/shared/contracts"
)

func handleTripStart(w http.ResponseWriter, r *http.Request) {
	var reqBody startTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	tripService, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		http.Error(w, "failed to connect to trip service", http.StatusInternalServerError)
		return
	}

	defer tripService.Close()

	tripStart, err := tripService.Client.CreateTrip(r.Context(), reqBody.toProto())
	if err != nil {
		http.Error(w, "failed to start trip: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := contracts.APIResponse{Data: tripStart}

	writeJSON(w, http.StatusCreated, response)
}

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

	tripService, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		http.Error(w, "failed to connect to trip service", http.StatusInternalServerError)
		return
	}

	defer tripService.Close()

	// resp, err := http.Post("http://trip-service:8083/preview", "application/json", reader)
	// if err != nil {
	// 	log.Print(err)
	// 	return
	// }

	// defer resp.Body.Close()

	tripPreview, err := tripService.Client.PreviewTrip(r.Context(), reqBody.toProto())
	if err != nil {
		http.Error(w, "failed to preview trip: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := contracts.APIResponse{Data: tripPreview}

	writeJSON(w, http.StatusCreated, response)
}
