package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	tripTypes "ride-sharing/services/trip-service/pkg/types"
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

func (s *service) GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*tripTypes.OsrmAPIResponse, error) {
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

	var routeResp tripTypes.OsrmAPIResponse
	if err := json.Unmarshal(body, &routeResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal OSRM response: %w", err)
	}
	return &routeResp, nil
}

func (s *service) EstimatePackagesWithRoute(route *tripTypes.OsrmAPIResponse) []*domain.RideFareModel {
	baseFares := getBaseFares()
	estimatedFares := make([]*domain.RideFareModel, len(baseFares))
	for i, f := range baseFares {
		estimatedFares[i] = estimateFareRoute(f, route)
	}
	return estimatedFares
}

func (s *service) GenerateTripFares(ctx context.Context, rideFares []*domain.RideFareModel, userID string) ([]*domain.RideFareModel, error) {
	fares := make([]*domain.RideFareModel, len(rideFares))

	for i, f := range rideFares {
		id := primitive.NewObjectID()

		fare := &domain.RideFareModel{
			UserID:            userID,
			ID:                id,
			TotalPriceInCents: f.TotalPriceInCents,
			PackageSlug:       f.PackageSlug,
		}

		if err := s.repo.SaveRideFare(ctx, fare); err != nil {
			return nil, fmt.Errorf("failed to save trip fare: %w", err)
		}

		fares[i] = fare
	}
	return fares, nil
}

func estimateFareRoute(f *domain.RideFareModel, route *tripTypes.OsrmAPIResponse) *domain.RideFareModel {
	pricingCfg := tripTypes.DefaultPricingConfig()
	carPackagePrice := f.TotalPriceInCents

	distanceKm := route.Routes[0].Distance
	distanceInMinutes := route.Routes[0].Duration

	distanceFare := distanceKm * pricingCfg.PricePerUnitOfDistance
	timeFare := distanceInMinutes * pricingCfg.PricingPerMinute
	totalFare := carPackagePrice + distanceFare + timeFare

	return &domain.RideFareModel{
		TotalPriceInCents: totalFare,
		PackageSlug:       f.PackageSlug,
	}
}

func getBaseFares() []*domain.RideFareModel {
	return []*domain.RideFareModel{
		{
			PackageSlug:       "suv",
			TotalPriceInCents: 200,
		},
		{
			PackageSlug:       "sedan",
			TotalPriceInCents: 350,
		},
		{
			PackageSlug:       "van",
			TotalPriceInCents: 450,
		},
		{
			PackageSlug:       "luxury",
			TotalPriceInCents: 1000,
		},
	}
}
