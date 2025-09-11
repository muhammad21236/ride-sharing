package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/shared/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	repo domain.TripRepository
}

func NewService(repo domain.TripRepository) *service {
	return &service{repo: repo}
}

func (s *service) CreateTrip(ctx context.Context, fare *domain.RideFareModel) (*domain.TripModel, error) {
	t := &domain.TripModel{
		ID:       primitive.NewObjectID(),
		UserID:   fare.UserID,
		Status:   "pending",
		RideFare: fare,
	}
	return s.repo.CreateTrip(ctx, t)
}

func (s *service) GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*types.OsrmAPIResponse, error) {
	url := "http://router.project-osrm.org/route/v1/driving/" +
		fmt.Sprintf("%f,%f;", pickup.Longitude, pickup.Latitude) +
		fmt.Sprintf("%f,%f", destination.Longitude, destination.Latitude) +
		"?overview=full&geometries=geojson"

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get route from OSRM: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read OSRM response body: %w", err)
	}

	var routeResp types.OsrmAPIResponse
	if err := json.Unmarshal(body, &routeResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal OSRM response: %w", err)
	}
	return &routeResp, nil
}
