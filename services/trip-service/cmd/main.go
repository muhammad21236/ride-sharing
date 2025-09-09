package main

import (
	"context"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"time"
)

func main() {
	ctx := context.Background()
	inmemRepo := repository.NewInmemRepostory()
	svc := service.NewService(inmemRepo)
	fare := &domain.RideFareModel{
		UserID: "user123",
	}
	t, err := svc.CreateTrip(ctx, fare)
	if err != nil {
		panic(err)
	}
	println("Created trip with ID:", t.ID.Hex())

	for {
		time.Sleep(1 * time.Second)
	}

}
